package find

import (
	"github.com/spf13/cobra"
	"srn/internal/shared"
)

var FileCmd = shared.NewCommand(
	"file <path> <terms...>",
	"Search for files by name",
	cobra.MinimumNArgs(2),
	func(cmd *cobra.Command, args []string) {
		path := args[0]
		terms := args[1:]
		FindAndProcess(path, terms, "f", "📄 Searching for files with '%s' in %s\n", "⚠️  Delete matched files? (y/N): ", false)
	},
)

var FileDeleteCmd = shared.NewCommand(
	"delete <path> <terms...>",
	"Delete files by name",
	cobra.MinimumNArgs(2),
	func(cmd *cobra.Command, args []string) {
		path := args[0]
		terms := args[1:]
		FindAndProcess(path, terms, "f", "📄 Searching for files with '%s' in %s\n", "⚠️  Delete matched files? (y/N): ", true)
	},
)

func init() {
	FileCmd.AddCommand(FileDeleteCmd)
}
