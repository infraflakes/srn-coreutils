package cmd

import (
	"github.com/spf13/cobra"

	todo "github.com/infraflakes/srn-todo/cmd"
)

func init() {
	rootCmd.AddCommand(todo.RootCmd)
}
