package cmd

import (
	"github.com/spf13/cobra"
	"srn/internal/nix"
)

func init() {
	rootCmd.AddCommand(nixCmd)

	nixCmd.AddCommand(nix.HomeCmd)
	nixCmd.AddCommand(nix.SysCmd)
	nix.BasicNixCmds(nixCmd)
}

var nixCmd = &cobra.Command{
	Use:   "nix",
	Short: "Nix related commands",
	Long:  `A collection of commands to manage Nix and NixOS systems.`,
}
