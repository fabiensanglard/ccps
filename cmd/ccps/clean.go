package main

import (
	"os"

	"github.com/spf13/cobra"
)

func rm(cmd *cobra.Command, dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		cmd.Printf("Cannot delete '%s' folder: %v\n", dir, err)
		os.Exit(1)
	}
	cmd.Println("rm -fr", dir)
}

func clean(cmd *cobra.Command, args []string) {
	rm(cmd, ".tmp")
	rm(cmd, "out")
}
