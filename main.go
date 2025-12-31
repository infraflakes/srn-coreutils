package main

import (
	"fmt"
	"os"

	"github.com/infraflakes/srn-coreutils/cmd"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && (os.Args[1] == "-v" || os.Args[1] == "--version") {
		fmt.Printf("serein coreutils version: %s\n", version)
		return
	}
	cmd.Execute()
}
