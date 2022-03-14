package z80

import (
	"bufio"
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

	obj := objiffy(linked)
	rom := pad(obj)

	return rom
}

func pad(input string) string {

	fi, err := os.Stat(input)
	if err != nil {
		println(fmt.Sprintf("Error stating '%s': %v", input, err))
		os.Exit(1)
	}
	// get the size
	if fi.Size() > board.Z80.Size {
		fmt.Printf("Z-80 RO< is too big (%d bytes) max=%d bytes", fi.Size(), board.Z80.Size)
	}

	cmd := fmt.Sprintf("dd if=/dev/zero of=%d bs=1 count=1 seek=65536", input)

	if verbose {
		println(cmd)
	}
	return ""
}

func objiffy(input string) string {
	output := objectDir + "z80.rom"
	cmd := fmt.Sprintf("%s --input-target=ihex --output-target=binary %s %s", objcopy, input, output)

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
		log.Fatal(err)
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
	linkerScript = append(linkerScript, "")

	for _, rel := range rels {
		linkerScript = append(linkerScript, fmt.Sprintf("%s", rel))
	}

	linkerScript = append(linkerScript, "")
	linkerScript = append(linkerScript, "-e")
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

	cmd := "sdldz80 -nf " + output

	if verbose {
		println(cmd)
	}

	return nil, output
}
