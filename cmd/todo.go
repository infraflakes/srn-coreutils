package cmd

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

	"srn/internal/todo"
)

var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Manage your todo list",
	Long:  `A terminal-based todo list manager with contexts, priorities, and more.`,
	Run: func(cmd *cobra.Command, args []string) {
		p := tea.NewProgram(todo.Initialize(), tea.WithAltScreen())

		if _, err := p.Run(); err != nil {
			fmt.Printf("Error running todo program: %v", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(todoCmd)
}
