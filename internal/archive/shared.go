package archive

import (
	"os"
	"path/filepath"
	"github.com/infraflakes/srn-libs/exec"
	"github.com/infraflakes/srn-libs/utils"
)

func ExpandTargets(targets []string) []string {
	for i, target := range targets {
		info, err := os.Stat(target)
		if err == nil && info.IsDir() {
			targets[i] = filepath.Join(target, "/*")
		}
	}
	return targets
}

func BuildArchiveCommand(archiveName string, targets []string, password string) {
	fileExt := filepath.Ext(archiveName)
	cmdArgs := []string{"a"}

	if password != "" {
		cmdArgs = append(cmdArgs, "-p"+password)
	}

	if fileExt == ".zip" {
		cmdArgs = append(cmdArgs, "-tzip")
	}

	cmdArgs = append(cmdArgs, archiveName)
	cmdArgs = append(cmdArgs, targets...)

	utils.CheckErr(exec.ExecuteCommand("7z", cmdArgs...))
}

func BuildExtractCommand(target string, password string) {
	cmdArgs := []string{"x"}
	if password != "" {
		cmdArgs = append(cmdArgs, "-p"+password)
	}
	cmdArgs = append(cmdArgs, target)
	utils.CheckErr(exec.ExecuteCommand("7z", cmdArgs...))
}
