package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"serein/internal/box"
)

var useDocker bool
var usePodman bool

func init() {
	rootCmd.AddCommand(boxCmd)

	boxCmd.PersistentFlags().BoolVar(&useDocker, "docker", false, "Force use of docker")
	boxCmd.PersistentFlags().BoolVar(&usePodman, "podman", false, "Force use of podman")

	box.StandaloneFlags(boxCmd)
	boxCmd.AddCommand(box.ContainerImagesCmd)
	boxCmd.AddCommand(box.ContainerIosCmd)
	boxCmd.AddCommand(box.ShellCmd)
	boxCmd.AddCommand(box.SilentCmd)
}

var boxCmd = &cobra.Command{
	Use:   "box",
	Short: "Manage containers with podman or docker aliases",
	Long:  `Manage containers with podman or docker aliases.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if useDocker && usePodman {
			fmt.Println("Error: --docker and --podman flags cannot be used simultaneously")
			os.Exit(1)
		}
		box.SetContainerEngine(useDocker, usePodman)
	},
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}
