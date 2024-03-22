package main

import (
	"encoding/binary"
	"fmt"
	"os"

	"github.com/fabiensanglard/ccps/boards"
	"github.com/fabiensanglard/ccps/oki"
	"github.com/fabiensanglard/ccps/sites"
	"github.com/spf13/cobra"
)

func dumpSFX(cmd *cobra.Command, args []string) {
	if targetBoard == "" {
		cmd.Println("No board target provided. Aborting")
		os.Exit(1)
	}

	dumpFolder := "dump/sfx/"
	if err := os.RemoveAll(dumpFolder); err != nil {
		cmd.Println("Unable to delete dump folder", dumpFolder)
		os.Exit(1)
	}
	if err := os.MkdirAll(dumpFolder, 0777); err != nil {
		cmd.Printf("Unable to create SFX dump folder '%s' : '%s'\n", dumpFolder, err.Error())
		os.Exit(1)
	}

	board := boards.Get(targetBoard)

	// Read the full ROM
	rom := make([]byte, board.Oki.Size)
	romCursor := 0
	for _, inRom := range board.Oki.Roms {
		path := sites.OutDir + inRom.Filename
		bytes, err := os.ReadFile(path)
		if err != nil {
			panic(fmt.Sprintf("Unable to read '%s' (%s)", path, err))
		}
		copy(rom[romCursor:romCursor+inRom.Size], bytes)
		romCursor += inRom.Size
	}

	// Parse and output wav
	okiRom := oki.OpenOKI(rom)
	if verbose {
		cmd.Printf("Found %d phrases\n", len(okiRom.Phrases))
	}
	for i, phrase := range okiRom.Phrases {
		if verbose {
			cmd.Printf("   Entry %d is %d bytes\n", i, len(phrase))
		}
		pcm := oki.ADPCMToPCM(phrase)
		writeWav(cmd, fmt.Sprintf("%s%03d.wav", dumpFolder, i), pcm)
	}
}

func writeWav(cmd *cobra.Command, path string, pcm []int16) {
	const wavHeaderSize = 0
	wav := make([]byte, wavHeaderSize+len(pcm)*2)

	// Master RIFF chunk
	copy(wav[0:4], "RIFF")
	binary.LittleEndian.PutUint32(wav[4:8], uint32(4+(8+16)+(8+len(pcm))))
	copy(wav[8:12], "WAVE")

	// fmt Chunk
	copy(wav[12:16], "fmt ")
	binary.LittleEndian.PutUint32(wav[16:20], 16)     // Subchunk1Size
	binary.LittleEndian.PutUint16(wav[20:22], 1)      // Format code PCM=0x0001
	binary.LittleEndian.PutUint16(wav[22:24], 1)      // Num Channels
	binary.LittleEndian.PutUint32(wav[24:28], 7575)   // Sampling rate
	binary.LittleEndian.PutUint32(wav[28:32], 7575*2) // Byte rate SampleRate * NumChannels * BitsPerSample/8
	binary.LittleEndian.PutUint16(wav[32:34], 2)      // Block Align (NumChannels * BitsPerSample/8)
	binary.LittleEndian.PutUint16(wav[34:36], 16)     // Bits per sample

	// data Chunk
	copy(wav[36:40], "data")
	payload := toByteArray(pcm)
	binary.LittleEndian.PutUint32(wav[40:44], uint32(len(payload))) // Sampling rate
	copy(wav[44:], payload)

	err := os.WriteFile(path, wav, 0644)
	if err != nil {
		cmd.Printf("Unable to write wav file at '%s'\n", path)
		os.Exit(1)
	}
}

func toByteArray(pcm []int16) []uint8 {
	data := make([]byte, 2*len(pcm))
	for i := 0; i < len(pcm); i += 1 {
		sample := pcm[i]
		data[i*2] = byte(sample & 0xFF)
		data[i*2+1] = byte(sample >> 8)
	}
	return data
}
