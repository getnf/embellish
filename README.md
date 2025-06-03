# Note

This application has been completely ported from Go to GJS (GNOME JavaScript), with focus shifted exclusively to the GUI implementation. The CLI/TUI components have been intentionally removed for the following strategic reasons:

- **[getnf](https://github.com/getnf/getnf.git)** provides a mature, feature-complete CLI/TUI solution that accomplishes the same objectives as our previous Go implementation, but with a significantly smaller bundle size and more streamlined architecture
- The Go codebase was becoming increasingly complex and difficult to maintain as new features were added, compounded by limited documentation and learning resources for GTK development in Go
- Development experience and maintainability concerns motivated the transition to a more suitable technology stack

This new GJS-based version offers enhanced polish, improved visual design, and introduces several new features while maintaining all core functionality. For users who prefer command-line interfaces, we strongly recommend **getnf** as the preferred CLI/TUI solution.

![Embellish Application Icon](/data/icons/io.github.getnf.embellish.svg)

# Embellish

[![Please do not theme this app](https://stopthemingmy.app/badge.svg)](https://stopthemingmy.app)

*A modern, intuitive font management application for Nerd Fonts*

## Screenshots

### Main Interface
![Main application interface showcasing the clean, modern design](/data/screenshots/main-page.png)

### Font Search & Discovery
![Advanced font search functionality with real-time filtering](/data/screenshots/search-page.png)

## Key Features

**Font Management**
- Browse and discover all available Nerd Fonts with detailed information
- One-click download and automatic system installation
- Safe removal of installed fonts with confirmation dialogs
- Seamless updates for previously installed fonts

**User Experience**
- Live font preview with customizable sample text
- Comprehensive license information display for legal compliance
- Intelligent search functionality with multiple filter options
- Responsive, native GTK4/Libadwaita interface following GNOME HIG

**Technical Excellence**
- Modern GJS/GTK4 architecture for optimal performance
- Flatpak distribution ensuring consistent cross-platform experience
- Minimal system resource usage with efficient font caching

## Installation

### Recommended: Flathub Installation

Embellish is officially distributed through Flathub, ensuring secure, sandboxed execution and automatic updates:

[<img width="240" alt="Download on Flathub" src="https://flathub.org/api/badge?svg&locale=en"/>](https://flathub.org/apps/io.github.getnf.embellish)

**Installation Command:**
```bash
flatpak install flathub io.github.getnf.embellish
```

**Launch Application:**
```bash
flatpak run io.github.getnf.embellish
```

### System Requirements
- Linux distribution with Flatpak support
- GTK4 and Libadwaita runtime
- Network connectivity for font downloads

## Development & Contributing

We welcome contributions from developers of all experience levels! This project adheres to the **[GNOME Code of Conduct](https://conduct.gnome.org)** to ensure a welcoming, inclusive development environment.

### Getting Started
- Fork the repository and create feature branches
- Follow existing code style and architectural patterns
- Test thoroughly across different desktop environments
- Submit pull requests with clear descriptions of changes

### Development Setup
```bash
# Clone repository
git clone https://github.com/getnf/embellish.git
cd embellish

# Install development dependencies
# (specific instructions depend on your distribution)
```

## Roadmap
- [ ] Translate to other languages

## Support & Community

- **Issues**: Report bugs and request features via GitHub Issues
- **Documentation**: Comprehensive user guides and developer documentation
- **Community**: Join discussions in project forums and Matrix channels

## License & Legal

This project respects all font licensing requirements and provides transparent license information for each available font. Users are responsible for compliance with individual font licenses in their specific use cases.

---

*Embellish: Making beautiful typography accessible to everyone*
