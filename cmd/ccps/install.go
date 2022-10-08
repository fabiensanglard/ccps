package main

import (
	"os"
	"path"

	"github.com/fabiensanglard/ccps/sites"
	"github.com/spf13/cobra"
)

func install(cmd *cobra.Command, args []string) {

	if installDestDir == "" {
		cmd.Println("missing destination directory")
		os.Exit(1)
	}

	srcDir := sites.OutDir
	files, err := os.ReadDir(srcDir)
	if err != nil {
		cmd.PrintErr(err)
		os.Exit(1)
	}

	if verbose {
		cmd.Println("Installing:")
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		src := path.Join(srcDir, file.Name())
		dst := path.Join(installDestDir, file.Name())
		if verbose {
			cmd.Println("Moving image '", src, "' -> '", dst, "'")
		}
		if err := os.Rename(src, dst); err != nil {
			cmd.Printf("Unable to move '%s' to '%s': '%s'\n", src, dst, err.Error())
			os.Exit(1)
		}
	}
}
