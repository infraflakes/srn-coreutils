package cmd

import (
	"github.com/infraflakes/srn-libs/exec"
	"github.com/infraflakes/srn-libs/utils"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "srn",
	Short: "Serein Coreutils is a suite for your cli.",
	Long:  `Serein Coreutils is a suite that supercharge your cli workflow.`,
	Run: func(cmd *cobra.Command, args []string) {
		utils.CheckErr(cmd.Help())
	},
}

func Execute() {
	utils.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&exec.DryRun, "dry-run", false, "print the command that would be executed instead of executing it")
}
