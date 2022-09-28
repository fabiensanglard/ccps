package m68k

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fabiensanglard/ccps/boards"
	"github.com/fabiensanglard/ccps/sites"
)

const cc = "m68k-linux-gnu-gcc"
const as = "m68k-linux-gnu-as"

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
		panic(fmt.Sprintf("Could not find executable '%s'", bin))
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
}

var verbose bool
var board boards.Board

func Build(v bool, b *boards.Board) []byte {
	verbose = v
	board = *b
	checkTools()
	sites.EnsureM68kObjsDir()

	var objs []string

	err, asmed := assemble()
	if err != nil {
		panic(fmt.Sprintf("Assembling error %s", err.Error()))
	}
	objs = append(objs, asmed...)

	// Compile user provider source code
	err, cced := compile(sites.M68kSrcsDir)
	if err != nil {
		panic(fmt.Sprintf("Compiling error '%s'", err.Error()))
	}
	objs = append(objs, cced...)

	// Compile generated source code (GFX assets)
	err, cced = compile(sites.M68kGenDir)
	if err != nil {
		panic(fmt.Sprintf("Compiling error '%s'", err.Error()))
	}
	objs = append(objs, cced...)

	err, linked := link(objs)
	if err != nil {
		panic(fmt.Sprintf("Linking error '%s'", err.Error()))
	}

	//romPath := binarize(linked)

	rom, err := os.ReadFile(linked)
	if err != nil {
		panic(fmt.Sprintf("Cannot read generated m68k ROM: '%s'", err.Error()))
	}

	return rom
}

//go:embed cps1.lk
var linkerScript []byte

func link(objs []string) (error, string) {
	lkPath := sites.M68kObjsDir + "cps1.lk"
	err := os.WriteFile(lkPath, linkerScript, 0644)
	if err != nil {
		panic(fmt.Sprintf("Unable to write linker script '%s'", lkPath))
	}

	mapDir := sites.M68kObjsDir + "game.map"
	output := sites.M68kObjsDir + "game.a"
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

func compile(dir string) (error, []string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
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
		output := sites.M68kObjsDir + name + ext_obj
		input := dir + src.Name()
		cmd := fmt.Sprintf("%s -DCPSB_VERSION=%d -I%s -m68000 -nostdlib -c -O0 -o %s %s",
			cc,
			board.Cpsb,
			sites.M68kGenDir,
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
	files, err := os.ReadDir(sites.M68kSrcsDir)
	if err != nil {
		panic(fmt.Sprintf("Unable to read dir '%s'", sites.M68kSrcsDir))
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
		output := sites.M68kObjsDir + name + ext_obj
		cmd := fmt.Sprintf("%s -m68000 --register-prefix-optional -o %s %s",
			as,
			output,
			sites.M68kSrcsDir+src.Name())
		run(cmd)

		if verbose {
			fmt.Println(cmd)
		}

		outputs = append(outputs, output)
	}

	return nil, outputs
}
