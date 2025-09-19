# 🌿 Branch Switcher

A beautiful TUI (Terminal User Interface) for managing Git branches across multiple repositories, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) using the Elm Architecture pattern.

![Branch Switcher Demo](docs/demo.gif)

## ✨ Features

- **🏗️ Elm Architecture**: Pure functional Model/Update/View pattern with immutable state
- **🚀 Auto-select all projects** for quick batch operations
- **🎯 Two operation modes**:
  1. Switch to main and pull latest changes
  2. Switch to main, pull latest, and create new branch
- **🎨 Beautiful TUI** with colors, icons, and intuitive navigation
- **⌨️ Vim-style keyboard shortcuts** for power users
- **🔍 Smart project detection** - automatically finds all Git repositories
- **⚡ Fast parallel operations** across multiple repositories
- **🛡️ Safe operations** - automatically stashes changes before switching

## 🚀 Quick Start

### Prerequisites
- Git repositories in your working directory

### Installation

#### Option 1: One-liner install (Recommended)
```bash
curl -sSL https://raw.githubusercontent.com/hoangbits/branch-switcher/main/install.sh | bash
```
*Automatically creates `brs` shortcut for quick access*

#### Option 2: Go install (if you have Go)
```bash
go install github.com/hoangbits/branch-switcher@latest
```

#### Option 3: Download binary
Download the latest release from [GitHub Releases](https://github.com/hoangbits/branch-switcher/releases)

#### Option 4: Build from source
```bash
git clone https://github.com/hoangbits/branch-switcher.git
cd branch-switcher
go build -o branch-switcher
chmod +x branch-switcher
sudo mv branch-switcher /usr/local/bin/
```

## 📖 Usage

### Basic Usage
Run from any repository within your multi-repository folder:

```bash
branch-switcher
# or use the short version
brs
```

The tool will automatically:
1. Scan the parent directory for Git repositories (containing `.git` folders)
2. Present an interactive menu with all repositories auto-selected
3. Allow you to choose your operation and execute it across selected repositories

### Navigation & Controls

#### Action Selection
- `↑/↓` or `j/k`: Navigate between options
- `Enter`: Select action
- `q` or `Ctrl+C`: Quit

#### Repository Selection
- `↑/↓` or `j/k`: Navigate through repository list
- `Space`: Toggle individual repository selection
- `a`: Toggle all repositories (select/deselect all)
- `Enter`: Confirm selection and proceed
- `Esc`: Go back to action selection
- `q` or `Ctrl+C`: Quit

#### Branch Name Input
- Type characters to enter branch name
- `Backspace`: Delete characters
- `Enter`: Confirm and start processing
- `Esc`: Go back to project selection

## 🏗️ Architecture

This project implements the **Elm Architecture** pattern, making the codebase predictable and maintainable:

```go
// Model: Immutable application state
type model struct {
    projects    []project        // Available projects
    selected    map[int]bool     // Selection state
    cursor      int              // Current cursor position
    mode        mode             // Current UI mode
    action      int              // Selected action (0: main, 1: branch)
    branchName  string           // Input branch name
    processing  bool             // Processing state
    error       string           // Error message
}

// Update: Pure state transition function
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd)

// View: Render current state to string
func (m model) View() string
```

### State Flow
```
User Input → Message → Update Function → New Model → View Function → Render
     ↑                                                                    ↓
     └────────────── User sees updated UI ←─────────────────────────────┘
```

### Key Benefits
- **Predictable**: All state changes go through the update function
- **Debuggable**: Easy to trace state transitions
- **Testable**: Pure functions are easy to unit test
- **Maintainable**: Clear separation of concerns

## 🔧 Development

### Project Structure
```
branch-switcher/
├── main.go           # Main application with Elm architecture
├── go.mod           # Go module dependencies
├── README.md        # This file
├── LICENSE          # MIT License
├── docs/            # Documentation and assets
│   ├── demo.gif     # Demo animation
│   └── architecture.md
└── .github/
    └── workflows/
        └── release.yml  # Automated releases
```

### Building
```bash
# Development build
go build

# Production build with optimizations
go build -ldflags "-s -w" -o branch-switcher

# Cross-platform builds
GOOS=linux GOARCH=amd64 go build -o branch-switcher-linux
GOOS=darwin GOARCH=arm64 go build -o branch-switcher-darwin
GOOS=windows GOARCH=amd64 go build -o branch-switcher.exe
```

### Dependencies
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework with Elm architecture
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style and layout library

### Contributing
1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Make your changes following the Elm architecture pattern
4. Test thoroughly
5. Commit: `git commit -m 'feat: add amazing feature'`
6. Push: `git push origin feature/amazing-feature`
7. Open a Pull Request

## 🎯 Use Cases

### Development Team Workflows
- **Sprint Start**: Switch all projects to main and pull latest changes
- **Feature Development**: Create feature branches across multiple services
- **Release Preparation**: Ensure all projects are on latest main
- **Hotfix Deployment**: Quickly create hotfix branches

### DevOps & CI/CD
- **Environment Sync**: Ensure all services are on correct branches
- **Deployment Preparation**: Batch branch operations before releases
- **Repository Maintenance**: Clean up and sync multiple repositories

## 🚨 Safety Features

- **Automatic Stashing**: Uncommitted changes are safely stashed before branch switching
- **Error Handling**: Clear error messages for failed operations
- **Non-destructive**: Never force-deletes branches with uncommitted changes
- **Rollback Support**: Easy to recover from failed operations

## 📝 Examples

### Switch all projects to main
```bash
branch-switcher
# 1. Select "Switch to main and pull latest"
# 2. All projects are auto-selected
# 3. Press Enter to execute
```

### Create feature branch across services
```bash
branch-switcher
# 1. Select "Switch to main, pull latest, and create new branch"
# 2. Deselect projects you don't need (all auto-selected)
# 3. Enter branch name: "feature/new-payment-system"
# 4. Press Enter to execute
```

## 🐛 Troubleshooting

### Common Issues

**"No Git repositories found in parent directory"**
- Ensure you're running from within a directory that has sibling Git repositories
- Check that repositories contain `.git` folders

**"Failed to fetch" errors**
- Verify internet connection
- Check Git remote configurations
- Ensure proper authentication for repositories

**Branch creation fails**
- Branch name may already exist
- Check for naming conflicts
- Ensure you have write permissions

### Debug Mode
```bash
# Run with verbose output
DEBUG=1 branch-switcher
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Charm](https://charm.sh/) for the amazing Bubble Tea framework
- [Elm Language](https://elm-lang.org/) for the architectural inspiration
- The Go community for excellent tooling and libraries

## 📞 Support

- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/hoangbits/branch-switcher/issues)
- 💡 **Feature Requests**: [GitHub Discussions](https://github.com/hoangbits/branch-switcher/discussions)
- 📧 **Contact**: [hoangbits@gmail.com](mailto:hoangbits@gmail.com)

---

Made with ❤️ using the Elm Architecture in Go