package main

import (
	"ccps/sites"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func install(args []string) {
	fs := flag.NewFlagSet("install", flag.ExitOnError)
	v := fs.Bool("v", false, "Verbose mode")
	dir := fs.String("d", "", "Destination directory")
	if err := fs.Parse(args); err != nil {
		panic(fmt.Sprintf("Cmd parsing error '%s'", err))
	}

	outDir := *dir
	if len(outDir) == 0 {
		panic("Usage: ccps install -d DIRECTORY")
	}
	verbose := *v
	if !strings.HasSuffix(outDir, "/") {
		outDir = outDir + "/"
	}

	srcDir := sites.OutDir
	files, err := ioutil.ReadDir(srcDir)
	if err != nil {
		log.Fatal(err)
	}

	if verbose {
		println("Installing:")
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		src := srcDir + file.Name()
		dst := outDir + file.Name()
		if verbose {
			println("Moving image '", src, "' -> '", dst, "'")
		}
		err := os.Rename(src, dst)
		if err != nil {
			panic(fmt.Sprintf("Unable to move '%s' to '%s': '%s'", src, dst, err.Error()))
		}
	}
}
