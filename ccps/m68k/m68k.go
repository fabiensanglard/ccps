package m68k

import (
	"bytes"
	"ccps/boards"
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const cc = "m68k-linux-gnu-gcc"
const as = "m68k-linux-gnu-as"

//const objcopy = "m68k-linux-gnu-objcopy"

const SrcsPath = "cc/68000/"
const objectDir = ".tmp/" + SrcsPath

const ext_as = ".s"
const ext_obj = ".o"
const ext_c = ".c"

func run(c string) {
	args := strings.Split(c, " ")
	cmd := exec.Command(args[0], args[1:]...)

	var out bytes.Buffer
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		println(fmt.Sprintf("Error running cmd '%s'", c))
		println(out.String())
		os.Exit(1)
	}

}

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
	//checkExecutable(objcopy)

	err := os.MkdirAll(objectDir, os.ModePerm)
	if err != nil {
		fmt.Println("Unable to create dir", objectDir)
		os.Exit(1)
	}
}

var verbose bool
var board boards.Board

func Build(v bool, b *boards.Board) []byte {
	verbose = v
	board = *b
	checkTools()

	var objs []string

	err, asmed := assemble()
	if err != nil {
		println("Assembling error", err)
		os.Exit(1)
	}
	objs = append(objs, asmed...)

	err, cced := compile()
	if err != nil {
		println("Compiling error", err)
		os.Exit(1)
	}
	objs = append(objs, cced...)

	err, linked := link(objs)
	if err != nil {
		println("Linking error", err)
		os.Exit(1)
	}

	//romPath := binarize(linked)

	rom, err := os.ReadFile(linked)
	if err != nil {
		println("Cannot read generated m68k ROM", err)
		os.Exit(1)
	}

	return rom
}

//func binarize(input string) string {
//	// TODO Check rom size before padding it.
//	// Get the size.
//	fi, err := os.Stat(input)
//	if err != nil {
//		println(fmt.Sprintf("Error stating '%s': %v", input, err))
//		os.Exit(1)
//	}
//
//	// Make sure it is not too big.
//	if fi.Size() > board.M68k.Size {
//		fmt.Printf("68000 ROM is too big (%d bytes) max=%d bytes", fi.Size(), board.M68k.Size)
//	}
//
//	output := objectDir + "game.rom"
//	// TODO: Double check why we remove .data via -R
//	cmd := fmt.Sprintf("%s --gap-fill=0xFF --pad-to=%d -R .data --output-target=binary %s %s", objcopy, board.M68k.Size, input, output)
//	run(cmd)
//	return output
//}

//go:embed cps1.lk
var linkerScript []byte

func link(objs []string) (error, string) {
	lkPath := objectDir + "cps1.lk"
	err := os.WriteFile(lkPath, linkerScript, 0644)
	if err != nil {
		println(fmt.Sprintf("Unable to write linker script '%s'", lkPath))
		os.Exit(1)
	}

	mapDir := objectDir + "game.map"
	output := objectDir + "game.a"
	cmd := fmt.Sprintf("%s -Llib -m68000 -Wall -nostartfiles -nodefaultlibs -fno-builtin -fomit-frame-pointer -ffast-math -Wl,-Map,%s -Wl,--build-id=none -T %s -o %s",
		cc,
		mapDir,
		lkPath,
		output)
	for _, obj := range objs {
		cmd = cmd + " " + obj
	}
	run(cmd)

	if verbose {
		fmt.Println(cmd)
	}

	return nil, output
}

func compile() (error, []string) {
	files, err := ioutil.ReadDir(SrcsPath)
	if err != nil {
		log.Fatal(err)
	}

	var outputs []string
	for _, src := range files {
		if src.IsDir() {
			continue
		}

		if !strings.HasSuffix(src.Name(), ext_c) {
			continue
		}

		basename := filepath.Base(src.Name())
		name := strings.TrimSuffix(basename, filepath.Ext(basename))
		output := objectDir + name + ext_obj
		input := SrcsPath + src.Name()
		cmd := fmt.Sprintf("%s -m68000 -nostdlib -c -O0 -o %s %s",
			cc,
			output,
			input)
		run(cmd)

		if verbose {
			fmt.Println(cmd)
		}

		outputs = append(outputs, output)
	}
	return nil, outputs

}

func assemble() (error, []string) {
	files, err := ioutil.ReadDir(SrcsPath)
	if err != nil {
		println(fmt.Sprintf("Unable to read dir '%s'", SrcsPath))
		os.Exit(1)
	}

	//TODO make sure crt0.s is first so areas are properly sorted.
	var outputs []string
	for _, src := range files {
		if src.IsDir() {
			continue
		}

		if !strings.HasSuffix(src.Name(), ext_as) {
			continue
		}

		basename := filepath.Base(src.Name())
		name := strings.TrimSuffix(basename, filepath.Ext(basename))
		output := objectDir + name + ext_obj
		cmd := fmt.Sprintf("%s -m68000 --register-prefix-optional -o %s %s",
			as,
			output,
			SrcsPath+src.Name())
		run(cmd)

		if verbose {
			fmt.Println(cmd)
		}

		outputs = append(outputs, output)
	}

	return nil, outputs
}
