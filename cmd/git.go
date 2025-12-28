package cmd

import (
	"github.com/spf13/cobra"
	"srn/internal/git"
)

func init() {
	rootCmd.AddCommand(gitCmd)

	git.BasicGitCommands(gitCmd)
	gitCmd.AddCommand(git.TagCmd)
	gitCmd.AddCommand(git.BranchCmd)
	gitCmd.AddCommand(git.CommitCmd)
}

var gitCmd = &cobra.Command{
	Use:   "git",
	Short: "Git helper commands",
	Long:  `A set of helper commands to simplify common Git operations.`,
}
