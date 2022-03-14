package m68k

import (
	"ccps/boards"
	"fmt"
	"os"
	"os/exec"
)

const cc = "m68k-linux-gnu-gcc"
const as = "m68k-linux-gnu-as"
const objcopy = "m68k-linux-gnu-objcopy"
const crt0 = "crt0.s"

func checkExecutable(bin string) {
	path, err := exec.LookPath(bin)
	if err != nil {
		fmt.Println("Could not find ", bin)
		os.Exit(1)
	}
	if verbose {
		fmt.Println(fmt.Sprintf("Found '%s' -> '%s'", bin, path))
	}
}

func checkTools() {
	if verbose {
		println("M68000 tools:")
	}
	checkExecutable(as)
	checkExecutable(cc)
	checkExecutable(objcopy)
}

var verbose bool

func Build(v bool, dryRun bool, board *boards.Board) []string {
	verbose = v

	checkTools()
	return nil
}
