package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/url"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/cli/go-gh"
	"github.com/cli/go-gh/v2/pkg/api"
)

var (
	operation            string
	auth                 struct{ Login string }
	selectedRepositories []string
	confirmedRandomText  string
	actionLabels         = map[string]interface{}{
		"delete":  "deleted",
		"archive": "archived",
	}
	inProgressActionLabels = map[string]interface{}{
		"delete":  "Deleteing",
		"archive": "Archiving",
	}
)

func main() {
	checkTokenScope()

	client, err := api.DefaultRESTClient()
	if err != nil {
		fmt.Println("Error creating API client:", err)
		os.Exit(0)
		return
	}

	err = client.Get("user", &auth)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
		return
	}

	for {
		operation = selectOperation()

		if operation == "exit" {
			os.Exit(0)
		}
		repoOptions := fetchReposOptions(client)

		if len(repoOptions) == 0 {
			fmt.Println("No Repositories Found.")
			continue
		}

		selectedRepositories = selectRepositories(repoOptions)

		if len(selectedRepositories) == 0 {
			fmt.Println("No Repositories Found.")
			continue
		}

		fmt.Println("Selected Repositories: ")
		for _, repo := range selectedRepositories {
			fmt.Println(repo)
		}

		confirmText()

		isConfirmed := confirmAction()

		if isConfirmed {
			err = spinner.New().
				Title(fmt.Sprintf("%s repositories...", inProgressActionLabels[operation])).
				Action(func() {
					if operation == "delete" {
						deleteRepos(client, selectedRepositories)
					} else if operation == "archive" {
						archiveRepos(client, selectedRepositories)
					}
				}).
				Run()

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func checkTokenScope() {
	statusContent, _, err := gh.Exec("auth", "status")

	if err != nil {
		log.Fatal(err)
		os.Exit(0)
	}

	// Convert output to string and split into lines
	outputStr := string(statusContent.String())
	lines := strings.Split(outputStr, "\n")

	// Find the line that contains "Token scopes"
	var scopesLine string
	for _, line := range lines {
		if strings.Contains(line, "Token scopes") {
			scopesLine = line
			break
		}
	}

	if scopesLine == "" {
		fmt.Println("No Token scopes found in the output.")
		os.Exit(0)
	}

	// Check if "delete_repo" is in the scopes
	hasDeleteRepo := strings.Contains(scopesLine, "delete_repo")

	if !hasDeleteRepo {
		fmt.Println("You don't have 'delete_repo' scope in your gh token.")
		fmt.Println("Refresh your token by running `gh auth refresh -s delete_repo` in your cli.")
		os.Exit(0)
	}
}

func randomBase64String(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/float64(1.33333333333))))
	rand.Read(buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:l] // strip 1 extra character we get from odd length results
}

func selectOperation() string {
	var value string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose operation").
				Options(
					huh.NewOption("Bulk: Delete repositories", "delete"),
					huh.NewOption("Bulk: Archive repositories", "archive"),
					huh.NewOption("Exit", "exit"),
				).
				Value(&value),
		),
	)

	err := form.Run()

	if err != nil {
		os.Exit(0)
	}

	return value
}

func selectRepositories(repoOptions []huh.Option[string]) []string {
	var selections []string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select Repositories to Process").
				Options(repoOptions...).
				Filterable(false).
				Value(&selections).
				Height(8),
		),
	).
		WithTheme(huh.ThemeDracula())

	err := form.Run()

	if err != nil {
		os.Exit(0)
	}

	return selections
}

func confirmText() string {
	randomText := randomBase64String(8)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Confirm the text?").
				Description(randomText).
				Prompt("? ").
				Validate(func(s string) error {
					if s != randomText {
						return errors.New("sorry, confirmation text didn't match")
					}

					return nil
				}).
				Value(&confirmedRandomText),
		),
	)

	err := form.Run()

	if err != nil {
		os.Exit(0)
	}

	return confirmedRandomText
}

func confirmAction() bool {
	var isConfirmed bool

	repoLabel := "repo"
	if len(selectedRepositories) > 1 {
		repoLabel = "repos"
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewConfirm().
				Title(fmt.Sprintf("❌ %d %s will be %s. Want to proceed? ❌", len(selectedRepositories), repoLabel, actionLabels[operation])).
				Affirmative("Yes!").
				Negative("No.").
				Value(&isConfirmed),
		),
	)

	err := form.Run()

	if err != nil {
		os.Exit(0)
	}

	return isConfirmed
}

func fetchReposOptions(client *api.RESTClient) []huh.Option[string] {
	var searchQuery string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Search").
				Prompt("? ").
				Description("Empty query will return 100 repositories").
				Value(&searchQuery),
		),
	)

	ierr := form.
		Run()

	if ierr != nil {
		fmt.Println(ierr)
		os.Exit(0)
	}

	fmt.Println("Fetching repositories...")

	queryParams := url.QueryEscape(fmt.Sprintf("%s user:%s", searchQuery, auth.Login))
	queryParams = strings.Trim(queryParams, " ")

	var result map[string]interface{}
	err := client.Get("search/repositories?q="+queryParams, &result)

	if err != nil {
		log.Panic("Error fetching repositories:", err)
		os.Exit(0)
	}

	var opions []huh.Option[string]

	// Extract repository names from the search result
	if items, ok := result["items"].([]interface{}); ok {
		for _, item := range items {
			if repo, ok := item.(map[string]interface{}); ok {
				if name, ok := repo["name"].(string); ok {
					opions = append(opions, huh.NewOption(name, name))
				}
			}
		}
	}

	return opions
}

func deleteRepos(client *api.RESTClient, repos []string) {
	for _, repo := range repos {
		var resp string
		var err = client.Delete("repos/"+auth.Login+"/"+repo, &resp)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Deleted " + repo)
	}
}

func archiveRepos(client *api.RESTClient, repos []string) {
	body := map[string]interface{}{
		"archived": true,
	}
	jsonBody, err := json.Marshal(body)

	if err != nil {
		log.Fatalf("impossible to build request: %s", err)
	}

	for _, repo := range repos {
		var resp interface{}
		err := client.Patch("repos/"+auth.Login+"/"+repo, bytes.NewReader(jsonBody), &resp)

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Archived " + repo)
	}
}
