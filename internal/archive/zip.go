package archive

import (
	"fmt"
	"github.com/spf13/cobra"
	"srn/internal/shared"
)

func init() {
	ZipCmd.AddCommand(zipPasswordCmd)
}

var ZipCmd = shared.NewCommand(
	"zip [archive-name] [target-to-archive]",
	"Archive files with 7z",
	cobra.MinimumNArgs(2),
	func(cmd *cobra.Command, args []string) {
		archiveName := args[0]
		targets := ExpandTargets(args[1:])
		BuildArchiveCommand(archiveName, targets, "")
	},
)

var zipPasswordCmd = shared.NewCommand(
	"password [archive-name] [target-to-archive]",
	"Archive files with 7z and a password",
	cobra.MinimumNArgs(2),
	func(cmd *cobra.Command, args []string) {
		archiveName := args[0]
		targets := ExpandTargets(args[1:])
		password := shared.GetInput("Enter archive password: ")
		if password == "" {
			fmt.Println("Password cannot be empty.")
			return
		}
		BuildArchiveCommand(archiveName, targets, password)
	},
)
