# gh-ibulk

✨ An interactive GitHub (gh) CLI extension to delete multiple repos at once.

<a href="https://github.com/mahmudz/gh-ibulk/releases"><img src="https://img.shields.io/github/release/mahmudz/gh-ibulk.svg" alt="Latest Release"></a>

![Github Bulk Delete Demo](https://raw.githubusercontent.com/mahmudz/gh-ibulk/refs/heads/main/demo.gif "Github Bulk Delete Demo")

## 📦 Installing from Github CLI

1. Install the `gh` CLI - see the [installation](https://github.com/cli/cli#installation)

   _Installation requires a minimum version (2.0.0) of the GitHub CLI that supports extensions._

2. Install this extension:

   ```sh
   gh extension install mahmudz/gh-ibulk
   ```

## ⚙️ Installing Manually

> If you want to install this extension **manually**, follow these steps:

1. Clone the repo

   ```shell
   # git
   git clone https://github.com/mahmudz/gh-ibulk
   ```

   ```shell
   # GitHub CLI
   gh repo clone mahmudz/gh-ibulk
   ```

2. Cd into it

   ```bash
   cd gh-ibulk
   ```

3. Build it

   ```bash
   go build
   ```

4. Install it locally
   ```bash
   gh extension install .
   ```

## ⚡️ Usage

Run & select your operation

```sh
gh ibulk
```
