package cmd

import (
	zip "github.com/infraflakes/srn-zip/cmd"
)

func init() {
	zip.RootCmd.Use = "zip"
	rootCmd.AddCommand(zip.RootCmd)
}
