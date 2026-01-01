package cmd

import (
	todo "github.com/infraflakes/srn-todo/cmd"
)

func init() {
	rootCmd.AddCommand(todo.RootCmd)
}
