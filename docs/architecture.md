# Architecture Documentation

## Overview

Branch Switcher follows the **Elm Architecture** pattern, providing a predictable and maintainable approach to building interactive terminal applications.

## The Elm Architecture

The Elm Architecture consists of three main components:

### 1. Model (State)

The Model represents the entire application state at any given moment:

```go
type model struct {
    projects    []project        // List of available Elixir projects
    selected    map[int]bool     // Which projects are selected
    cursor      int              // Current cursor position in UI
    mode        mode             // Current UI mode (action select, project select, etc.)
    action      int              // Selected action (0: main, 1: branch)
    branchName  string           // User-entered branch name
    processing  bool             // Whether operations are in progress
    error       string           // Current error message to display
}
```

### 2. Update (State Transitions)

The Update function is the only way to modify the application state. It takes the current model and a message, then returns a new model and optional command:

```go
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
```

### 3. View (Rendering)

The View function takes the current model and renders it to a string that represents the UI:

```go
func (m model) View() string {
    switch m.mode {
    case modeSelectAction:
        return m.renderActionSelect()
    case modeSelectProjects:
        return m.renderProjectSelect()
    case modeEnterBranch:
        return m.renderBranchInput()
    }
}
```

## State Flow Diagram

```
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│  User Input │───▶│   Message   │───▶│   Update    │
└─────────────┘    └─────────────┘    └─────────────┘
                                              │
┌─────────────┐    ┌─────────────┐    ┌─────────────┐
│   Render    │◀───│    View     │◀───│ New Model   │
└─────────────┘    └─────────────┘    └─────────────┘
```

## Application Modes

The application has several distinct modes, each with its own UI and behavior:

### 1. `modeSelectAction`
- **Purpose**: Choose between "switch to main" or "create branch"
- **Navigation**: Up/Down arrows
- **Actions**: Enter to select, Q to quit

### 2. `modeSelectProjects`
- **Purpose**: Select which projects to operate on
- **Features**: All projects auto-selected by default
- **Navigation**: Up/Down arrows, Space to toggle
- **Actions**: Enter to continue, Esc to go back

### 3. `modeEnterBranch`
- **Purpose**: Get branch name from user
- **Input**: Character-by-character input
- **Actions**: Enter to confirm, Esc to go back

### 4. `modeProcessing`
- **Purpose**: Show progress while executing Git operations
- **Behavior**: Non-interactive, shows spinner/progress

## Message Types

The application uses typed messages for all state transitions:

```go
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
```

## Benefits of This Architecture

### 1. **Predictability**
- All state changes flow through the Update function
- No hidden state mutations
- Easy to reason about application behavior

### 2. **Debuggability**
- State transitions are explicit and traceable
- Each message represents a specific user action
- Easy to add logging or debugging

### 3. **Testability**
- Update function is pure (given same inputs, produces same outputs)
- View function is pure
- Easy to unit test individual components

### 4. **Maintainability**
- Clear separation between state, logic, and presentation
- Adding new features follows established patterns
- Refactoring is safer due to type system

## Error Handling

Errors are handled as part of the application state:

```go
// Errors become part of the model
type model struct {
    error string
    // ... other fields
}

// Errors are displayed in the view
func (m model) View() string {
    if m.error != "" {
        return errorStyle.Render("❌ " + m.error) + "\n\n"
    }
    // ... rest of view
}

// Errors clear automatically after display
m.error = "" // Clear after showing
```

## Asynchronous Operations

Long-running operations (like Git commands) are handled using Bubble Tea's command system:

```go
func (m model) processProjects(branchName string) tea.Cmd {
    return tea.Tick(tea.Millisecond*100, func(t tea.Time) tea.Msg {
        // Perform Git operations
        for id, selected := range m.selected {
            if selected {
                err := switchBranch(m.projects[id].path, branchName)
                if err != nil {
                    return msg{Type: msgError, Error: err}
                }
            }
        }
        return msg{Type: msgProcessComplete}
    })
}
```

## Style System

UI styling is handled through Lip Gloss:

```go
var (
    titleStyle = lipgloss.NewStyle().
        Background(lipgloss.Color("62")).
        Foreground(lipgloss.Color("230")).
        Padding(0, 1)

    selectedStyle = lipgloss.NewStyle().
        Foreground(lipgloss.Color("86"))
)
```

This ensures consistent theming and makes it easy to modify the visual appearance.

## Future Extensibility

The architecture makes it easy to add new features:

1. **New Modes**: Add new mode constants and update functions
2. **New Messages**: Add message types and handle them in Update
3. **New Views**: Create new rendering functions
4. **New Operations**: Add new Git operations following the command pattern

The Elm Architecture scales well and maintains clarity even as applications grow in complexity.