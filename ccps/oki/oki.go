package oki

import (
	"ccps/boards"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func Build(v bool, dryRun bool, board *boards.Board) string {
	verbose := v

	wavDir := "sfx/"
	files, err := ioutil.ReadDir(wavDir)
	if err != nil {
		println("Unable to open gfx dir", wavDir)
		os.Exit(1)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if !strings.HasSuffix(file.Name(), ".wav") {
			if verbose {
				println("Skipping non-wav '", file.Name(), "'")
			}
			continue
		}

		if verbose {
			println("Processing wav '", file.Name(), "'")
		}

		wavPath := wavDir + file.Name()
		wav, err := LoadWav(wavPath)
		if err != nil {
			println("Unable to open wav file", wavPath)
			os.Exit(1)
		}

		// sox test.wav -r 7575 -b 8 -c 1 outfile.wav
		// soxi sfx/moldova.wav
		// mpg321 -w moldova.wav moldova.mp3
		adpcm := toADPCM(wav.data)
		println(len(adpcm))
	}

	outDir := "out/"
	rom := make([]byte, board.Oki.Size)

	romPath := outDir + "gfx.rom"
	err = ioutil.WriteFile(romPath, rom, 0644)
	if err != nil {
		fmt.Println("Unable to write Oki rom to", romPath)
		os.Exit(1)
	}

	// Return everything rom path
	return romPath
}
