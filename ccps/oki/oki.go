package oki

import (
	"ccps/boards"
	"ccps/code"
	"ccps/sites"
	"io/ioutil"
	"os"
	"strings"
)

func Build(v bool, board *boards.Board) ([]byte, *code.Code) {
	verbose := v

	files, err := ioutil.ReadDir(sites.SfxSrcPath)
	if err != nil {
		if verbose {
			println("Unable to open sfx dir", sites.SfxSrcPath)
		}
		return nil, code.NewCode()
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

		wavPath := sites.SfxSrcPath + file.Name()
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

	// TODO: Generate Oki ID headers. Returning empty for now
	okiIDHeader := code.NewCode()

	return okiRom.genROM(board.Oki.Size), okiIDHeader
}
