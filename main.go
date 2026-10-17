package main

import (
	"os"

	"github.com/ctrlpad/daemon/cmd"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	os.Exit(run())
}

func run() int {
	cmd.Init(version, commit, date)

	err := cmd.Execute()
	if err != nil {
		return 1
	}

	return 0
}
