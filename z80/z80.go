package z80

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fabiensanglard/ccps/boards"
	"github.com/fabiensanglard/ccps/sites"
)

const cc = "sdcc"
const as = "sdasz80"
const linker = "sdldz80"
const objcopy = "objcopy"
const dd = "dd"

const ext_as = ".s"
const ext_rel = ".rel"
const ext_c = ".c"

var verbose bool
var board boards.Board

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

func Build(v bool, b *boards.Board) []byte {
	verbose = v
	board = *b

	checkTools()

	var rels []string

	err, asmed := assemble()
	if err != nil {
		panic(fmt.Sprintf("Assembling error '%s'", err.Error()))
	}
	rels = append(rels, asmed...)

	err, cced := compile()
	if err != nil {
		panic(fmt.Sprintf("Compiling error '%s'", err.Error()))
	}
	rels = append(rels, cced...)

	err, linked := link(rels)
	if err != nil {
		panic(fmt.Sprintf("Linking error '%s'", err.Error()))
	}

	obj := binarize(linked)
	romPath := pad(obj)

	rom, err := os.ReadFile(romPath)
	if err != nil {
		panic(fmt.Sprintf("Cannot read generated z80 ROM '%s'", err.Error()))
	}

	return rom
}

func pad(input string) string {

	// Get the size.
	fi, err := os.Stat(input)
	if err != nil {
		panic(fmt.Sprintf("Error stating '%s': %v", input, err))
	}

	// Make sure it is not too big.
	if fi.Size() > board.Z80.Size {
		panic(fmt.Sprintf("Z-80 ROM is too big (%d bytes) max=%d bytes", fi.Size(), board.Z80.Size))
	}

	cmd := fmt.Sprintf("dd if=/dev/zero of=%s bs=1 count=1 seek=65535", input)
	run(cmd)

	if verbose {
		println(cmd)
	}
	return input
}

func binarize(input string) string {
	output := sites.Z80ObjsDir + "z80.rom"
	cmd := fmt.Sprintf("%s --input-target=ihex --output-target=binary %s %s", objcopy, input, output)
	run(cmd)

	if verbose {
		println(cmd)
	}

	return output
}

func checkExecutable(bin string) {
	path, err := exec.LookPath(bin)
	if err != nil {
		panic(fmt.Sprintf("Could not find ", bin))
	}
	if verbose {
		fmt.Println(fmt.Sprintf("Found '%s' -> '%s'", bin, path))
	}
}

func checkTools() {
	if verbose {
		println("Z-80 tools:")
	}
	checkExecutable(as)
	checkExecutable(cc)
	checkExecutable(linker)
	checkExecutable(objcopy)
	checkExecutable(dd)

	if verbose {
		println("Creating folder", sites.Z80ObjsDir)
	}

	sites.EnsureZ80ObjsDir()
}

// cc/z80/*.as -> .tmp/cc/z80/*.rel
func assemble() (error, []string) {

	files, err := os.ReadDir(sites.Z80SrcsDir)
	if err != nil {
		panic(fmt.Sprintf("Unable to read dir '%s'", sites.Z80SrcsDir))
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
		output := sites.Z80ObjsDir + name + ext_rel
		cmd := fmt.Sprintf("%s -plogff -o %s %s",
			as,
			output,
			sites.Z80SrcsDir+src.Name())
		run(cmd)

		if verbose {
			fmt.Println(cmd)
		}

		outputs = append(outputs, output)
	}

	return nil, outputs
}

func compile() (error, []string) {
	files, err := os.ReadDir(sites.Z80SrcsDir)
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
		output := sites.Z80ObjsDir + name + ext_rel
		input := sites.Z80SrcsDir + src.Name()
		cmd := fmt.Sprintf("%s --compile-only -mz80 --data-loc 0xd000 --no-std-crt0 -o %s %s",
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

// sdldz80 -nf main.lk
/*

-e
*/
func link(rels []string) (error, string) {
	output := sites.Z80ObjsDir + "main.ihx"
	var linkerScript []string
	linkerScript = append(linkerScript, "-mjwx")
	linkerScript = append(linkerScript, fmt.Sprintf("-i %s", output))
	linkerScript = append(linkerScript, "-b _CODE = 0x0200")
	linkerScript = append(linkerScript, "-b _DATA = 0xd000")
	linkerScript = append(linkerScript, "-k /usr/share/sdcc/lib/z80")
	linkerScript = append(linkerScript, "-l z80")

	for _, rel := range rels {
		linkerScript = append(linkerScript, fmt.Sprintf("%s", rel))
	}

	linkerScript = append(linkerScript, "")

	// Write linker file
	mode := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	lkPath := sites.Z80ObjsDir + "main.lk"
	file, err := os.OpenFile(lkPath, mode, 0644)
	defer file.Close()
	if err != nil {
		panic(fmt.Sprintf("Cannot create linker script '%s'", lkPath))
	}
	datawriter := bufio.NewWriter(file)
	for _, data := range linkerScript {
		_, err = datawriter.WriteString(data + "\n")
		if err != nil {
			panic(fmt.Sprintf("Cannot write linker script '%s'", err))
		}
	}
	datawriter.Flush()

	cmd := "sdldz80 -nf " + lkPath
	run(cmd)

	if verbose {
		println(cmd)
	}

	return nil, output
}
