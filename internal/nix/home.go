package nix

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/infraflakes/srn-libs/cli"
)

func init() {
	HomeCmd.AddCommand(HomeBuildCmd)
	HomeCmd.AddCommand(HomeGenCmd)
	HomeCmd.AddCommand(HomeGenDeleteCmd)
}

var HomeCmd = cli.NewCommand(
	"home",
	"Manage home-manager",
	cobra.NoArgs,
	func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
)

var HomeBuildCmd = cli.NewCommand(
	"build [path/to/flake]",
	"Build a home-manager configuration",
	cobra.ExactArgs(1),
	func(cmd *cobra.Command, args []string) {
		flakePath := args[0]
		runNixCommand("home-manager", "switch", "--flake", flakePath)
	},
)

var HomeGenCmd = cli.NewCommand(
	"gen",
	"List home-manager generations",
	cobra.NoArgs,
	func(cmd *cobra.Command, args []string) {
		runNixCommand("home-manager", "generations")
	},
)

var HomeGenDeleteCmd = cli.NewCommand(
	"delete [numbers...]",
	"Delete home-manager generations",
	cobra.MinimumNArgs(1),
	func(cmd *cobra.Command, args []string) {
		generations, err := parseGenerations(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cmdArgs := append([]string{"remove-generations"}, generations...)
		runNixCommand("home-manager", cmdArgs...)
	},
)
