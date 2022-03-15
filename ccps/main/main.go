package main

import (
	"fmt"
	"os"
)

const cmdBuild = "build"
const cmdHelloWorld = "hw"
const cmdInstall = "install"
const cmdClean = "clean"

var allCmds = []string{cmdBuild, cmdInstall, cmdHelloWorld, cmdClean}

func rm(dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		println(fmt.Sprintf("Cannot delete '%s' folder: %v", dir, err))
		os.Exit(1)
	}
	println("rm -fr", dir)
}
func clean() {
	rm(".tmp")
	rm("out")
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(fmt.Sprintf("Error: Expected subcommands %v", allCmds))
		os.Exit(1)
	}
	cmd := os.Args[1]
	args := os.Args[1:]

	if cmd == cmdBuild {
		build(args)
	} else if cmd == cmdHelloWorld {
		helloWorld(args)
	} else if cmd == cmdInstall {
		install(args)
	} else if cmd == cmdClean {
		clean()
	} else {
		println(fmt.Sprintf("Usage: ccps %v", allCmds))
	}
}
