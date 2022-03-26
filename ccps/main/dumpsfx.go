package main

import (
	"ccps/boards"
	"ccps/oki"
	"encoding/binary"
	"flag"
	"fmt"
	"os"
)

func dumpSFX(args []string) {
	fs := flag.NewFlagSet("sfx", flag.ContinueOnError)
	verbose := fs.Bool("v", false, "Verbose mode")
	boardName := fs.String("b", "", "Target board")

	if err := fs.Parse(args); err != nil {
		println(fmt.Sprintf("Cmd parsing error '%s'", err))
		os.Exit(1)
	}

	if *boardName == "" {
		println("No board target provided. Aborting")
		os.Exit(1)
	}

	dumpFolder := "dump/sfx/"
	err := os.RemoveAll(dumpFolder)
	//if err != nil {
	//	println("Unable to delete dump folder ", dumpFolder)
	//	os.Exit(1)
	//}
	err = os.MkdirAll(dumpFolder, 0777)
	if err != nil {
		println("Unable to create dump folder ", dumpFolder, ":", err.Error())
		os.Exit(1)
	}

	board := boards.Get(*boardName)

	// Read the full ROM
	rom := make([]byte, board.Oki.Size)
	romCursor := 0
	for _, inRom := range board.Oki.Roms {
		path := "out/" + inRom.Filename
		bytes, err := os.ReadFile(path)
		if err != nil {
			println(fmt.Sprintf("Unable to read '%s' (%s)", path, err))
			os.Exit(1)
		}
		copy(rom[romCursor:romCursor+inRom.Size], bytes)
		romCursor += inRom.Size
	}

	// Parse and output wav
	okiRom := oki.OpenOKI(rom)
	if *verbose {
		println(fmt.Sprintf("Found %d phrases", len(okiRom.Phrases)))
	}
	for i, phrase := range okiRom.Phrases {
		if *verbose {
			println(fmt.Sprintf("   Entry %d is %d bytes", i, len(phrase)))
		}
		pcm := oki.ADPCMToPCM(phrase)
		writeWav(fmt.Sprintf("%s%03d.wav", dumpFolder, i), pcm)
	}
}

func writeWav(path string, pcm []byte) {
	const wavHeaderSize = 0
	wav := make([]byte, wavHeaderSize+len(pcm))

	// Master RIFF chunk
	copy(wav[0:4], "RIFF")
	binary.LittleEndian.PutUint32(wav[4:8], uint32(4+(8+16)+(8+len(pcm))))
	copy(wav[8:12], "WAVE")

	// fmt Chunk
	copy(wav[12:16], "fmt ")
	binary.LittleEndian.PutUint32(wav[16:20], 16)   // Subchunk1Size
	binary.LittleEndian.PutUint16(wav[20:22], 1)    // Format code PCM=0x0001
	binary.LittleEndian.PutUint16(wav[22:24], 1)    // Num Channels
	binary.LittleEndian.PutUint32(wav[24:28], 7575) // Sampling rate
	binary.LittleEndian.PutUint32(wav[28:32], 7575) // Byte rate SampleRate * NumChannels * BitsPerSample/8
	binary.LittleEndian.PutUint16(wav[32:34], 1)    // Block Align (NumChannels * BitsPerSample/8)
	binary.LittleEndian.PutUint16(wav[34:36], 8)    // Bits per sample

	// data Chunk
	copy(wav[36:40], "data")
	binary.LittleEndian.PutUint32(wav[40:44], uint32(len(pcm))) // Sampling rate
	copy(wav[44:], pcm)

	err := os.WriteFile(path, wav, 0644)
	if err != nil {
		println("Unable to write", path)
		os.Exit(1)
	}
}
