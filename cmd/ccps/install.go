package main

import (
	"os"
	"strings"

	"github.com/fabiensanglard/ccps/sites"
	"github.com/spf13/cobra"
)

func install(cmd *cobra.Command, args []string) {

	if installDestDir == "" {
		cmd.Println("missing destination directory")
		os.Exit(1)
	}
	if !strings.HasSuffix(installDestDir, "/") {
		installDestDir = installDestDir + "/"
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

		src := srcDir + file.Name()
		dst := installDestDir + file.Name()
		if verbose {
			cmd.Println("Moving image '", src, "' -> '", dst, "'")
		}
		err := os.Rename(src, dst)
		if err != nil {
			cmd.Printf("Unable to move '%s' to '%s': '%s'\n", src, dst, err.Error())
			os.Exit(1)
		}
	}
}
