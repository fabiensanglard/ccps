package main

import (
	"flag"
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
		//println(fmt.Sprintf("Cmd parsing error '%s'", err))
		os.Exit(1)
	}

	outDir := *dir
	if len(outDir) == 0 {
		println("Usage: ccps install -d dstDirectory")
		os.Exit(1)
	}
	verbose := *v
	if !strings.HasSuffix(outDir, "/") {
		outDir = outDir + "/"
	}

	srcDir := "out/"
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
			println("Unable to move", src, "to", dst, ":", err.Error())
			os.Exit(1)
		}
	}
}
