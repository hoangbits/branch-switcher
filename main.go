package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1).
		MarginBottom(1)

	selectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("86"))

	unselectedStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241"))

	helpStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("241")).
		MarginTop(1)

	errorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color("196"))
)

// Messages (Elm-style)
type msgType int

const (
	msgToggleProject msgType = iota
	msgToggleAll
	msgConfirm
	msgSetMode
	msgSetBranchName
	msgProcessComplete
	msgError
)

type msg struct {
	Type      msgType
	ProjectID int
	Mode      mode
	Branch    string
	Error     error
}

type mode int

const (
	modeSelectAction mode = iota
	modeSelectProjects
	modeEnterBranch
	modeProcessing
)

// Model (Elm-style)
type model struct {
	projects    []project
	selected    map[int]bool
	cursor      int
	mode        mode
	action      int // 0: switch to main, 1: create branch
	branchName  string
	processing  bool
	error       string
}

type project struct {
	name string
	path string
}

// Update function (Elm-style)
func (m model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.KeyMsg:
		switch m.mode {
		case modeSelectAction:
			return m.updateActionSelect(msg)
		case modeSelectProjects:
			return m.updateProjectSelect(msg)
		case modeEnterBranch:
			return m.updateBranchInput(msg)
		}
	}
	return m, nil
}

func (m model) updateActionSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < 1 {
			m.cursor++
		}
	case "enter":
		m.action = m.cursor
		m.mode = modeSelectProjects
		m.cursor = 0
		// Auto-select all projects
		for i := range m.projects {
			m.selected[i] = true
		}
	}
	return m, nil
}

func (m model) updateProjectSelect(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.mode = modeSelectAction
		m.cursor = 0
		m.selected = make(map[int]bool)
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down", "j":
		if m.cursor < len(m.projects)-1 {
			m.cursor++
		}
	case " ":
		m.selected[m.cursor] = !m.selected[m.cursor]
	case "a":
		allSelected := len(m.selected) == len(m.projects)
		m.selected = make(map[int]bool)
		if !allSelected {
			for i := range m.projects {
				m.selected[i] = true
			}
		}
	case "enter":
		if len(m.selected) == 0 {
			m.error = "No projects selected"
			return m, nil
		}

		if m.action == 1 { // Create branch
			m.mode = modeEnterBranch
			m.branchName = ""
		} else { // Switch to main
			return m, m.processProjects("")
		}
	}
	return m, nil
}

func (m model) updateBranchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "q", "ctrl+c":
		return m, tea.Quit
	case "esc":
		m.mode = modeSelectProjects
	case "enter":
		if m.branchName == "" {
			m.error = "Branch name cannot be empty"
			return m, nil
		}
		return m, m.processProjects(m.branchName)
	case "backspace":
		if len(m.branchName) > 0 {
			m.branchName = m.branchName[:len(m.branchName)-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.branchName += msg.String()
		}
	}
	return m, nil
}

func (m model) processProjects(branchName string) tea.Cmd {
	return tea.Tick(tea.Millisecond*100, func(t tea.Time) tea.Msg {
		for id, selected := range m.selected {
			if !selected {
				continue
			}

			project := m.projects[id]
			if err := switchBranch(project.path, branchName); err != nil {
				return msg{Type: msgError, Error: err}
			}
		}
		return msg{Type: msgProcessComplete}
	})
}

func switchBranch(projectPath, branchName string) error {
	// Stash changes
	exec.Command("git", "-C", projectPath, "stash").Run()

	// Fetch and switch to main
	if err := exec.Command("git", "-C", projectPath, "fetch", "origin").Run(); err != nil {
		return fmt.Errorf("failed to fetch: %v", err)
	}

	exec.Command("git", "-C", projectPath, "branch", "-D", "main").Run()

	if err := exec.Command("git", "-C", projectPath, "checkout", "--track", "origin/main").Run(); err != nil {
		return fmt.Errorf("failed to checkout main: %v", err)
	}

	if err := exec.Command("git", "-C", projectPath, "pull", "origin", "main").Run(); err != nil {
		return fmt.Errorf("failed to pull: %v", err)
	}

	// Create new branch if specified
	if branchName != "" {
		if err := exec.Command("git", "-C", projectPath, "checkout", "-b", branchName).Run(); err != nil {
			return fmt.Errorf("failed to create branch: %v", err)
		}
	}

	return nil
}

// View function (Elm-style)
func (m model) View() string {
	var b strings.Builder

	// Title
	b.WriteString(titleStyle.Render("🌿 Branch Switcher"))
	b.WriteString("\n\n")

	// Error display
	if m.error != "" {
		b.WriteString(errorStyle.Render("❌ " + m.error))
		b.WriteString("\n\n")
		m.error = "" // Clear error after showing
	}

	switch m.mode {
	case modeSelectAction:
		b.WriteString(m.renderActionSelect())
	case modeSelectProjects:
		b.WriteString(m.renderProjectSelect())
	case modeEnterBranch:
		b.WriteString(m.renderBranchInput())
	case modeProcessing:
		b.WriteString("🔄 Processing projects...")
	}

	return b.String()
}

func (m model) renderActionSelect() string {
	var b strings.Builder

	b.WriteString("What would you like to do?\n\n")

	actions := []string{
		"Switch to main and pull latest",
		"Switch to main, pull latest, and create new branch",
	}

	for i, action := range actions {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
			action = selectedStyle.Render(action)
		} else {
			action = unselectedStyle.Render(action)
		}
		b.WriteString(fmt.Sprintf("%s %s\n", cursor, action))
	}

	b.WriteString(helpStyle.Render("\nUse ↑/↓ to navigate, enter to select, q to quit"))

	return b.String()
}

func (m model) renderProjectSelect() string {
	var b strings.Builder

	actionText := "switch to main"
	if m.action == 1 {
		actionText = "create new branch"
	}

	b.WriteString(fmt.Sprintf("Select projects to %s (all auto-selected):\n\n", actionText))

	for i, project := range m.projects {
		cursor := " "
		if i == m.cursor {
			cursor = ">"
		}

		checkbox := "[ ]"
		style := unselectedStyle
		if m.selected[i] {
			checkbox = "[✓]"
			style = selectedStyle
		}

		line := fmt.Sprintf("%s %s %s", cursor, checkbox, project.name)
		b.WriteString(style.Render(line) + "\n")
	}

	selectedCount := len(m.selected)
	b.WriteString(fmt.Sprintf("\nSelected: %d/%d", selectedCount, len(m.projects)))

	b.WriteString(helpStyle.Render("\nSpace to toggle, 'a' for all, enter to continue, esc to go back"))

	return b.String()
}

func (m model) renderBranchInput() string {
	var b strings.Builder

	b.WriteString("Enter branch name:\n\n")
	b.WriteString(fmt.Sprintf("> %s_\n\n", m.branchName))

	b.WriteString(helpStyle.Render("Type branch name, enter to continue, esc to go back"))

	return b.String()
}

// Init function (Elm-style)
func (m model) Init() tea.Cmd {
	return nil
}

func findProjects() []project {
	var projects []project

	// Get parent directory
	cwd, _ := os.Getwd()
	parentDir := filepath.Dir(cwd)

	// Find directories with .git folders (Git repositories)
	dirs, _ := os.ReadDir(parentDir)
	for _, dir := range dirs {
		if dir.IsDir() {
			gitPath := filepath.Join(parentDir, dir.Name(), ".git")
			if stat, err := os.Stat(gitPath); err == nil && stat.IsDir() {
				projects = append(projects, project{
					name: dir.Name(),
					path: filepath.Join(parentDir, dir.Name()),
				})
			}
		}
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].name < projects[j].name
	})

	return projects
}

func main() {
	projects := findProjects()

	if len(projects) == 0 {
		fmt.Println("No Git repositories found in parent directory")
		os.Exit(1)
	}

	initialModel := model{
		projects: projects,
		selected: make(map[int]bool),
		mode:     modeSelectAction,
	}

	p := tea.NewProgram(initialModel, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}