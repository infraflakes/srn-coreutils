package todo

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Configuration and persistence

func (m *Model) LoadConfig() {
	configDir := filepath.Dir(m.ConfigFilePath)
	if err := os.MkdirAll(configDir, 0755); err != nil { // owner can read, write, and execute
		fmt.Println("Error creating config directory:", err)
		m.CreateDefaultConfig()
		return
	}

	data, err := os.ReadFile(m.ConfigFilePath)
	if err != nil {
		m.CreateDefaultConfig()
		return
	}
	// Creating an anonymous struct right there inside the function to match the structure of your JSON data.
	var config struct {
		Tasks    []Task   `json:"tasks"`
		NextID   int      `json:"next_id"`
		Contexts []string `json:"contexts"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		m.CreateDefaultConfig()
		return
	}

	m.Tasks = config.Tasks
	m.NextID = config.NextID
	m.Contexts = config.Contexts

	if m.NextID == 0 { // Fallback to ensure IDs don't get messed up if the NextID value is somehow corrupted or missing, prevents duplicate IDs
		maxID := 0
		for _, task := range m.Tasks {
			if task.ID > maxID {
				maxID = task.ID
			}
		}
		m.NextID = maxID + 1
	}
}

// Creating an anonymous struct right there inside the function to match the structure of the JSON data.
// To void the need to define a separate Config type globally.
func (m *Model) SaveConfig() {
	configDir := filepath.Dir(m.ConfigFilePath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		fmt.Println("Error creating config directory:", err)
		return
	}
	config := struct {
		Tasks    []Task   `json:"tasks"`
		NextID   int      `json:"next_id"`
		Contexts []string `json:"contexts"`
	}{
		Tasks:    m.Tasks,
		NextID:   m.NextID,
		Contexts: m.Contexts,
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling config:", err)
		return
	}

	if err := os.WriteFile(m.ConfigFilePath, data, 0644); err != nil {
		fmt.Println("Error saving config file:", err)
	}
}

// Default behavior
func (m *Model) CreateDefaultConfig() {
	m.Tasks = []Task{
		{ID: 1, Task: "Welcome to your todo app!", Checked: false, Context: "Getting Started"},
		{ID: 2, Task: "Press 'a' to add a new task", Checked: false, Context: "Getting Started"},
		{ID: 3, Task: "Press space to toggle completion", Checked: true, Context: "Getting Started"},
		{ID: 4, Task: "Use arrow keys to navigate", Checked: false, Context: "Getting Started"},
		{ID: 5, Task: "Press '?' to see more keybindings", Checked: false, Context: "Getting Started"},
	}
	m.Contexts = []string{"Getting Started"}
	m.NextID = 6
}
