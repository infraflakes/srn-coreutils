package find

import (
	"fmt"
	"os"

	"srn/internal/shared"
)

func DeletePath(path string, isDir bool) {
	var err error
	if isDir {
		err = os.RemoveAll(path)
	} else {
		err = os.Remove(path)
	}
	if err != nil {
		fmt.Printf("Failed to delete %s: %v\n", path, err)
	} else {
		fmt.Printf("Deleted: %s\n", path)
	}
}

func FindAndProcess(path string, terms []string, findType string, searchMessage string, deletePrompt string, del bool) {
	isDir := findType == "d"
	for _, term := range terms {
		fmt.Printf(searchMessage, term, path)
		matches, err := shared.RunCommand("find", path, "-type", findType, "-name", fmt.Sprintf("*%s*", term))
		if err != nil {
			fmt.Printf("Error finding paths: %v\n", err)
			continue
		}
		for _, match := range matches {
			fmt.Println(match)
		}

		if del && len(matches) > 0 && shared.Confirm(deletePrompt) {
			for _, match := range matches {
				DeletePath(match, isDir)
			}
		}
	}
}

func DeleteGrepMatches(basePath, term string) {
	matches, err := shared.RunCommand("grep", "-rlE", term, basePath)
	if err != nil {
		fmt.Printf("Error finding files to delete: %v\n", err)
		return
	}
	for _, match := range matches {
		DeletePath(match, false)
	}
}
