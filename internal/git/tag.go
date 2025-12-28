package git

import (
	"github.com/spf13/cobra"
	"srn/internal/shared"
)

func init() {
	TagCmd.AddCommand(tagCreateCmd)
	TagCmd.AddCommand(tagDeleteCmd)
	TagCmd.AddCommand(tagWipeCmd)

	tagDeleteCmd.AddCommand(tagDeleteLocalCmd)
	tagDeleteCmd.AddCommand(tagDeleteRemoteCmd)
}

var TagCmd = shared.NewCommand(
	"tag",
	"Git tag commands",
	cobra.NoArgs,
	func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
)

var tagCreateCmd = shared.NewCommand(
	"create [SHA] [name] [message]",
	"Alias for git tag -a <name> -m \"<msg>\" <SHA>",
	cobra.ExactArgs(3),
	func(cmd *cobra.Command, args []string) {
		runGitCommand("tag", "-a", args[1], "-m", args[2], args[0])
	},
)

var tagDeleteCmd = shared.NewCommand(
	"delete",
	"Delete a local or remote tag",
	cobra.NoArgs,
	func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
)

var tagDeleteLocalCmd = shared.NewCommand(
	"local [name]",
	"Alias for git tag -d <name>",
	cobra.ExactArgs(1),
	func(cmd *cobra.Command, args []string) {
		runGitCommand("tag", "-d", args[0])
	},
)

var tagDeleteRemoteCmd = shared.NewCommand(
	"remote [name]",
	"Alias for git push origin --delete <name>",
	cobra.ExactArgs(1),
	func(cmd *cobra.Command, args []string) {
		runGitCommand("push", "origin", "--delete", args[0])
	},
)

var tagWipeCmd = shared.NewCommand(
	"wipe [name]",
	"Alias for git tag -d <name> && git push origin --delete <name>",
	cobra.ExactArgs(1),
	func(cmd *cobra.Command, args []string) {
		runGitCommand("tag", "-d", args[0])
		runGitCommand("push", "origin", "--delete", args[0])
	},
)
