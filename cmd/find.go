package cmd

import (
	find "github.com/infraflakes/srn-find/cmd"
)

func init() {
	rootCmd.AddCommand(find.RootCmd)
}
