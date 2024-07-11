<p align="center">
  <img src="https://raw.githubusercontent.com/getnf/embellish/main/io.github.getnf.embellish.svg" width="100" />
</p>
<p align="center">
    <h1 align="center" style="color: f5c211ff">Embellish</h1>
</p>
<p align="center">
    <em>Install Nerd Fonts</em>
    <br />
    <strong><em>This README file is a work in progress</em></strong>
</p>
<p align="center">
	<img src="https://img.shields.io/github/license/getnf/embellish?style=flat&color=0080ff" alt="license">
	<img src="https://img.shields.io/github/last-commit/getnf/embellish?style=flat&logo=git&logoColor=white&color=0080ff" alt="last-commit">
	<img src="https://img.shields.io/github/languages/top/getnf/embellish?style=flat&color=0080ff" alt="repo-top-language">
	<img src="https://img.shields.io/github/languages/count/getnf/embellish?style=flat&color=0080ff" alt="repo-language-count">
</p>
<p align="center">
		<em>Developed with the software and tools below.</em>
</p>
<p align="center">
	<img src="https://img.shields.io/badge/Go-1.22.3-%2301ADD8?logo=go" alt="Go">
    <img src="https://img.shields.io/badge/GTK-4-%233584E4?style=flat&logo=gtk&label=GTK&color=%233584E4" alt="Gtk4">
    <img src="https://img.shields.io/badge/libadwaita-1.4-%231C70D5?logo=gnome" alt="libadwaita-1.4">
</p>
<hr>

## 🔗 Quick Links

> - [📍 Overview](#-overview)
> - [📦 Features](#-features)
> - [🚀 Getting Started](#-getting-started)
> - [🤝 Contributing](#-contributing)
> - [📄 License](#-license)
> - [👏 Acknowledgments](#-acknowledgments)

---

## 📍 Overview

An application written in go, helps you install Nerd Fonts.

The app works as a CLI, TUI and GUI written with GTK4 and libadwita 1.4.

- TUI Install
  <img src="https://raw.githubusercontent.com/getnf/embellish/main/Screenshots/tui-install-prompt.png" width="500" />
- TUI uninstall
  <img src="https://raw.githubusercontent.com/getnf/embellish/main/Screenshots/tui-uninstall-prompt.png" width="500" />
- GUI main page
    <img src="https://raw.githubusercontent.com/getnf/embellish/main/Screenshots/main-page.png" width="500" />
- GUI search page
    <img src="https://raw.githubusercontent.com/getnf/embellish/main/Screenshots/search-page.png" width="500" />  

---

## 📦 Features

- List all available Nerd Fonts
- Download and install a Font
- Uninstall an already installed Font
- Update all installed fonts
- Bulk install/uninstall (cli/tui only)
- Fuzzy search (tui/gui)
- Font name suggestion (if typed wrong) (cli)
- Separate binaries for the GUI and CLI/TUI versions

---

## 🚀 Getting Started

***Requirements***

Ensure you have the following dependencies installed on your system:

* **go**: `version 1.22.3` 
* **gcc**
  
For the GUI build: 
* **gtk4-devel**
* **libadwaita-devel**
* **gobject-introspection-devel**
* **gdk-pixbuf2-devel**
* **graphene-devel**
* **at-devel**
* **atk-devel**
* **pango-devel**
* **gtk4-devel**
* **gtk3-devel**


### ⚙️ Installation

1. Clone the embellish repository:

```sh
git clone https://github.com/getnf/embellish
```

2. Change to the project directory:

```sh
cd embellish
```

3. Install the dependencies:

for GUI binary

```sh
go build -v -tags gui -o embelish-gui
```

for tui/cli binary

```sh
go build -v -tags terminal -o embelish-tui
```

you can also use `-ldflags "-s -w"` to output smaller size binaries with any debug info.

---

## 🤝 Contributing

Contributions are welcome! Here are several ways you can contribute:

- **[Submit Pull Requests](https://github.com/getnf/embellish/blob/main/CONTRIBUTING.md)**: Review open PRs, and submit your own PRs.
- **[Report Issues](https://github.com/getnf/embellish/issues)**: Submit bugs found or log feature requests for Embellish and provide feedback.

<details closed>
    <summary>Contributing Guidelines</summary>

1. **Fork the Repository**: Start by forking the project repository to your GitHub account.
2. **Clone Locally**: Clone the forked repository to your local machine using a Git client.
   ```sh
   git clone https://github.com/getnf/embellish
   ```
3. **Create a New Branch**: Always work on a new branch, giving it a descriptive name.
   ```sh
   git checkout -b new-feature-x
   ```
4. **Make Your Changes**: Develop and test your changes locally.
5. **Commit Your Changes**: Commit with a clear message describing your updates.
   ```sh
   git commit -m '[FEAT] Implemented new feature x.'
   ```
6. **Push to GitHub**: Push the changes to your forked repository.
   ```sh
   git push origin new-feature-x
   ```
7. **Submit a Pull Request**: Create a PR against the original project repository. Clearly describe the changes and their motivations.

</details>

---

## 📄 License

This project is protected under the [GPL-v3](https://choosealicense.com/licenses) License. For more details, refer to the [LICENSE](https://github.com/getnf/embellish/blob/main/LICENSE) file.

---

## 👏 Acknowledgments

- [diamondburned](https://github.com/diamondburned): for [gotk4](https://github.com/diamondburned/gotk4) and [gotk4-adwaita](https://github.com/diamondburned/gotk4-adwaita) go bindings for gtk4 and libadwiata
- [Ryanoasis](https://github.com/ryanoasis): for the amazing [Nerd Fonts](https://github.com/ryanoasis/nerd-fonts)

