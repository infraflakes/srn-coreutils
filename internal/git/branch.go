package git

import (
	"github.com/spf13/cobra"
	"srn/internal/shared"
)

func init() {
	BranchCmd.AddCommand(gitBranchListCmd)
	BranchCmd.AddCommand(gitBranchCreateCmd)
	BranchCmd.AddCommand(gitBranchSwitchCmd)
	BranchCmd.AddCommand(gitBranchDeleteLocalCmd)
	BranchCmd.AddCommand(gitBranchDeleteRemoteCmd)
}

var BranchCmd = shared.NewCommand(
	"branch",
	"Git branch commands",
	cobra.NoArgs,
	func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
)

var gitBranchListCmd = shared.NewCommand(
	"list",
	"Alias for git branch -avv",
	cobra.NoArgs,
	func(cmd *cobra.Command, args []string) {
		runGitCommand("branch", "-avv")
	},
)

var gitBranchCreateCmd = shared.NewCommand(
	"create [name]",
	"Alias for git switch -c <name>",
	cobra.ExactArgs(1),
	func(cmd *cobra.Command, args []string) {
		runGitCommand("switch", "-c", args[0])
	},
)

var gitBranchSwitchCmd = shared.NewCommand(
	"switch [name]",
	"Alias for git switch <name>",
	cobra.ExactArgs(1),
	func(cmd *cobra.Command, args []string) {
		runGitCommand("switch", args[0])
	},
)

var gitBranchDeleteLocalCmd = shared.NewCommand(
	"local [name]",
	"Alias for git branch -D <name>",
	cobra.ExactArgs(1),
	func(cmd *cobra.Command, args []string) {
		runGitCommand("branch", "-D", args[0])
	},
)

var gitBranchDeleteRemoteCmd = shared.NewCommand(
	"remote [name]",
	"Alias for git push origin --delete <name>",
	cobra.ExactArgs(1),
	func(cmd *cobra.Command, args []string) {
		runGitCommand("push", "origin", "--delete", args[0])
	},
)
