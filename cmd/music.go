package cmd

import (
	music "github.com/infraflakes/srn-music/cmd"
)

func init() {
	music.RootCmd.Use = "music"
	rootCmd.AddCommand(music.RootCmd)
}
