package nix

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"srn/internal/shared"
)

func init() {
	SysCmd.AddCommand(SysBuildCmd)
	SysCmd.AddCommand(SysGenCmd)
	SysCmd.AddCommand(SysGenDeleteCmd)
}

var SysCmd = shared.NewCommand(
	"sys",
	"Manage NixOS system",
	cobra.NoArgs,
	func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
)

var SysBuildCmd = shared.NewCommand(
	"build [path/to/flake]",
	"Build a NixOS system",
	cobra.ExactArgs(1),
	func(cmd *cobra.Command, args []string) {
		flakePath := args[0]
		runNixCommand("sudo", "nixos-rebuild", "switch", "--impure", "--flake", flakePath)
	},
)

var SysGenCmd = shared.NewCommand(
	"gen",
	"List system generations",
	cobra.NoArgs,
	func(cmd *cobra.Command, args []string) {
		runNixCommand("sudo", "nix-env", "--list-generations", "--profile", "/nix/var/nix/profiles/system")
	},
)

var SysGenDeleteCmd = shared.NewCommand(
	"delete [numbers...]",
	"Delete system generations",
	cobra.MinimumNArgs(1),
	func(cmd *cobra.Command, args []string) {
		generations, err := parseGenerations(args)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		cmdArgs := append([]string{"nix-env", "--profile", "/nix/var/nix/profiles/system", "--delete-generations"}, generations...)
		runNixCommand("sudo", cmdArgs...)
	},
)
