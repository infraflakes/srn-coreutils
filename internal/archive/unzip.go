package archive

import (
	"github.com/spf13/cobra"
	"srn/internal/shared"
)

func init() {
	UnzipCmd.AddCommand(unzipPasswordCmd)
}

var UnzipCmd = shared.NewCommand(
	"unzip [target-to-unarchive]",
	"Extract archives with 7z",
	cobra.ExactArgs(1),
	func(cmd *cobra.Command, args []string) {
		BuildExtractCommand(args[0], "")
	},
)

var unzipPasswordCmd = shared.NewCommand(
	"password [password] [target-to-unarchive]",
	"Extract archives with a password",
	cobra.ExactArgs(2),
	func(cmd *cobra.Command, args []string) {
		BuildExtractCommand(args[1], args[0])
	},
)
