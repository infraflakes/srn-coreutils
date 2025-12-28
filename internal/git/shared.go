package git

import (
	"srn/internal/shared"
)

func runGitCommand(args ...string) {
	shared.CheckErr(shared.ExecuteCommand("git", args...))
}
