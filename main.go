package main

import (
	"os"

	"github.com/spf13/cobra-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
