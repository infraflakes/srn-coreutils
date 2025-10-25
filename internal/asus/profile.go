package asus

import (
	"github.com/spf13/cobra"
	"serein/internal/shared"
)

func init() {
	AsusProfileCmd.AddCommand(AsusProfileStatusCmd)
	AsusProfileCmd.AddCommand(AsusProfileSetCmd)
	AsusProfileCmd.AddCommand(AsusProfileListCmd)
}

var AsusProfileCmd = shared.NewCommand(
	"profile",
	"Manage asusctl profiles",
	cobra.NoArgs,
	func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
)

var AsusProfileStatusCmd = shared.NewCommand(
	"status",
	"Shows current asusctl profile",
	cobra.NoArgs,
	func(cmd *cobra.Command, args []string) {
		runAsusCommand("asusctl", "profile", "-p")
	},
)

var AsusProfileSetCmd = shared.NewCommand(
	"set [profile]",
	"Set asusctl profiles",
	cobra.ExactArgs(1),
	func(cmd *cobra.Command, args []string) {
		asusProfile := args[0]
		runAsusCommand("asusctl", "profile", "-P", asusProfile)
	},
)

var AsusProfileListCmd = shared.NewCommand(
	"list",
	"List asusctl profiles",
	cobra.NoArgs,
	func(cmd *cobra.Command, args []string) {
		runAsusCommand("asusctl", "profile", "-l")
	},
)
