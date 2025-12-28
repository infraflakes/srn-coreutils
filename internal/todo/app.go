package todo

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
)

// Initialize creates a new model
func Initialize(configFilePath string) Model {
	var finalPath string
	if configFilePath != "" {
		finalPath, _ = filepath.Abs(configFilePath)
	} else {
		homeDir, _ := os.UserHomeDir()
		// New default path
		finalPath = filepath.Join(homeDir, ".cache", "srn-todo", "note.json")
	}

	ti := textinput.New()
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 70

	dateInputs := make([]textinput.Model, 3)
	for i := range dateInputs {
		dateInputs[i] = textinput.New()
		dateInputs[i].Focus()
		dateInputs[i].CharLimit = 4
		dateInputs[i].Width = 10
	}

	m := Model{
		TextInput:      ti,
		DateInputs:     dateInputs,
		KeyMap:         DefaultKeyMap(),
		Help:           help.New(),
		ConfigFilePath: finalPath, // Changed from ConfigPath
		MaxHistory:     50,
		ViewMode:       NormalView,
	}

	m.LoadConfig()
	m.UpdateContexts()

	return m
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}
