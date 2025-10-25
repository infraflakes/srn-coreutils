package asus

import (
	"serein/internal/shared"
)

func runAsusCommand(command string, args ...string) {
	shared.CheckErr(shared.ExecuteCommand(command, args...))
}
