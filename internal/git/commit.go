package git

import (
	"fmt"
	"github.com/spf13/cobra"
	"srn/internal/shared"
)

func init() {
	CommitCmd.AddCommand(CommitPushCmd)
	CommitCmd.AddCommand(CommitListCmd)
	CommitCmd.AddCommand(CommitUndoCmd)
	CommitCmd.AddCommand(commitDeleteCmd)
	CommitCmd.AddCommand(CommitChangesCmd)
	CommitCmd.AddCommand(CommitCompareCmd)

	commitDeleteCmd.AddCommand(commitDeleteStageCmd)
	commitDeleteCmd.AddCommand(commitDeleteUnstageCmd)
}

var CommitCmd = shared.NewCommand(
	"commit [message]",
	"Alias for git commit -m \"<commit msg>\"",
	cobra.ExactArgs(1),
	func(cmd *cobra.Command, args []string) {
		runGitCommand("commit", "-m", args[0])
	},
)

var CommitPushCmd = shared.NewCommand(
	"push <branch> [force]",
	"Alias for git push origin <branch>",
	cobra.RangeArgs(1, 2),
	func(cmd *cobra.Command, args []string) {
		branchName := args[0]
		isForce := false

		if len(args) == 2 {
			if args[1] == "force" {
				isForce = true
			} else {
				fmt.Println("Error: invalid second argument. Did you mean 'force'?")
				return
			}
		}

		gitArgs := []string{"push", "origin", branchName}
		if isForce {
			gitArgs = append(gitArgs, "--force")
		}
		runGitCommand(gitArgs...)
	},
)

var CommitListCmd = shared.NewCommand(
	"list",
	"Alias for git log --oneline",
	cobra.NoArgs,
	func(cmd *cobra.Command, args []string) {
		runGitCommand("log", "--oneline")
	},
)

var CommitUndoCmd = shared.NewCommand(
	"undo [SHA]",
	"Alias for git revert <SHA>",
	cobra.ExactArgs(1),
	func(cmd *cobra.Command, args []string) {
		runGitCommand("revert", args[0])
	},
)

// commitDeleteCmd is the new parent for delete subcommands
var commitDeleteCmd = shared.NewCommand(
	"delete",
	"Delete a commit using --soft or --mixed reset",
	cobra.NoArgs,
	func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
)

var commitDeleteStageCmd = shared.NewCommand(
	"stage [number]",
	"Alias for git reset --soft HEAD~<numb> (keeps changes staged)",
	cobra.ExactArgs(1),
	func(cmd *cobra.Command, args []string) {
		runGitCommand("reset", "--soft", "HEAD~"+args[0])
	},
)

var commitDeleteUnstageCmd = shared.NewCommand(
	"unstage [number]",
	"Alias for git reset --mixed HEAD~<numb> (keeps changes in working dir)",
	cobra.ExactArgs(1),
	func(cmd *cobra.Command, args []string) {
		runGitCommand("reset", "--mixed", "HEAD~"+args[0])
	},
)

var CommitChangesCmd = shared.NewCommand(
	"changes [SHA]",
	"Alias for git show <SHA>",
	cobra.ExactArgs(1),
	func(cmd *cobra.Command, args []string) {
		runGitCommand("show", args[0])
	},
)

var CommitCompareCmd = shared.NewCommand(
	"compare [SHA1] [SHA2]",
	"Alias for git diff <SHA1> <SHA2>",
	cobra.ExactArgs(2),
	func(cmd *cobra.Command, args []string) {
		runGitCommand("diff", args[0], args[1])
	},
)
