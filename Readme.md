# 🚀 gh-ibulk

✨ An interactive GitHub (gh) CLI extension to delete multiple repos at once.

<a href="https://github.com/mahmudz/gh-ibulk/releases"><img src="https://img.shields.io/github/release/mahmudz/gh-ibulk.svg" alt="Latest Release"></a>

![Github Bulk Delete Demo](https://raw.githubusercontent.com/mahmudz/gh-ibulk/refs/heads/main/demo.gif "Github Bulk Delete Demo")

## 🧩 Features
- 🗑️ Delete repos
- 🗃️ Archive repos

## 📦 Installing from Github CLI

1. Install the `gh` CLI - see the [installation](https://github.com/cli/cli#installation)

   _Installation requires a minimum version (2.0.0) of the GitHub CLI that supports extensions._

2. Install this extension:

   ```sh
   gh extension install mahmudz/gh-ibulk
   ```

## ⚙️ Installing Manually

> If you want to install this extension **manually**, follow these steps:

1. You must have go `go>=1.23.1` installed on your machine
2. Clone the repo

   ```shell
   # git
   git clone https://github.com/mahmudz/gh-ibulk
   ```

   ```shell
   # GitHub CLI
   gh repo clone mahmudz/gh-ibulk
   ```

3. Cd into it

   ```bash
   cd gh-ibulk
   ```

4. Build it

   ```bash
   go build
   ```

5. Install it locally
   ```bash
   gh extension install .
   ```

## ⚡️ Usage

Run & select your operation

```sh
gh ibulk
```

## 💬 How to offer feedback or make a feature request
If you would like to offer feedback or make a feature request, please create a new discussion [here](https://github.com/mahmudz/gh-ibulk/discussions/new/choose).

## 🤝🏼 Contribution
You're welcome to contribute.

## 👨🏻‍💻 Author
Mahmudul Hasan [hello@mahmudz.dev](hello@mahmudz.dev)
