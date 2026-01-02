package cmd

import (
	cd "github.com/infraflakes/srn-cd/cmd"
)

func init() {
	cd.RootCmd.Use = "cd"
	rootCmd.AddCommand(cd.RootCmd)
}
