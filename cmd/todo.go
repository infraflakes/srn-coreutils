package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"srn/internal/todo"
)

var todoCmd = &cobra.Command{
	Use:   "todo [path/to/note.json]",
	Short: "Manage your todo list",
	Long:  `A terminal-based todo list manager with contexts, priorities, and more.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var configPath string
		if len(args) > 0 {
			configPath = args[0]
		}
		p := tea.NewProgram(todo.Initialize(configPath), tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running todo program: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(todoCmd)
}
