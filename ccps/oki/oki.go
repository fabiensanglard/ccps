package oki

import (
	"ccps/boards"
	"io/ioutil"
	"os"
	"strings"
)

func Build(v bool, board *boards.Board) []byte {
	verbose := v

	wavDir := "sfx/"
	files, err := ioutil.ReadDir(wavDir)
	if err != nil {
		println("Unable to open gfx dir", wavDir)
		os.Exit(1)
	}

	okiRom := NewOkiROM()
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
		adpcm := PCMtoADPCM(wav.data)
		okiRom.AddPhrase(adpcm)
	}

	return okiRom.genROM(board.Oki.Size)
}
