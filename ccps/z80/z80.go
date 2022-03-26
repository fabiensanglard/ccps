package z80

import (
	"bufio"
	"bytes"
	"ccps/boards"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const cc = "sdcc"
const as = "sdasz80"
const linker = "sdldz80"
const objcopy = "objcopy"
const dd = "dd"

const srcsPath = "cc/z80/"
const objectDir = ".tmp/" + srcsPath

const ext_as = ".s"
const ext_rel = ".rel"
const ext_c = ".c"

var verbose bool
var dryRun bool
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

func Build(v bool, dr bool, b *boards.Board) string {
	verbose = v
	dryRun = dr
	board = *b

	checkTools()

	var rels []string

	err, asmed := assemble()
	if err != nil {
		println("Assembling error", err)
		os.Exit(1)
	}
	rels = append(rels, asmed...)

	err, cced := compile()
	if err != nil {
		println("Compiling error", err)
		os.Exit(1)
	}
	rels = append(rels, cced...)

	err, linked := link(rels)
	if err != nil {
		println("Linking error", err)
		os.Exit(1)
	}

	obj := binarize(linked)
	rom := pad(obj)

	return rom
}

func pad(input string) string {

	// Get the size.
	fi, err := os.Stat(input)
	if err != nil {
		println(fmt.Sprintf("Error stating '%s': %v", input, err))
		os.Exit(1)
	}

	// Make sure it is not too big.
	if fi.Size() > board.Z80.Size {
		fmt.Printf("Z-80 ROM is too big (%d bytes) max=%d bytes", fi.Size(), board.Z80.Size)
		os.Exit(1)
	}

	cmd := fmt.Sprintf("dd if=/dev/zero of=%s bs=1 count=1 seek=65535", input)
	run(cmd)

	if verbose {
		println(cmd)
	}
	return input
}

func binarize(input string) string {
	output := objectDir + "z80.rom"
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
		fmt.Println("Could not find ", bin)
		os.Exit(1)
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
		println("Creating folder", objectDir)
	}
	err := os.MkdirAll(objectDir, os.ModePerm)
	if err != nil {
		fmt.Println("Unable to create dir", objectDir)
		os.Exit(1)
	}
}

// cc/z80/*.as -> .tmp/cc/z80/*.rel
func assemble() (error, []string) {

	files, err := ioutil.ReadDir(srcsPath)
	if err != nil {
		println(fmt.Sprintf("Unable to read dir '%s'", srcsPath))
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
		output := objectDir + name + ext_rel
		cmd := fmt.Sprintf("%s -plogff -o %s %s",
			as,
			output,
			srcsPath+src.Name())
		run(cmd)

		if verbose {
			fmt.Println(cmd)
		}

		outputs = append(outputs, output)
	}

	return nil, outputs
}

func compile() (error, []string) {
	files, err := ioutil.ReadDir(srcsPath)
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
		output := objectDir + name + ext_rel
		input := srcsPath + src.Name()
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
	output := objectDir + "main.ihx"
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
	lkPath := objectDir + "main.lk"
	file, err := os.OpenFile(lkPath, mode, 0644)
	defer file.Close()
	if err != nil {
		println(fmt.Sprintf("Cannot create linker script '%s'", lkPath))
		os.Exit(1)
	}
	datawriter := bufio.NewWriter(file)
	for _, data := range linkerScript {
		_, err = datawriter.WriteString(data + "\n")
		if err != nil {
			println(fmt.Sprintf("Cannot write linker script '%s'", err))
			os.Exit(1)
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
