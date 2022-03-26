package oki

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
)

type WaveFmt struct {
	Subchunk1ID   string // Must be "fmt "
	Subchunk1Size uint32 // 16 for PCM
	AudioFormat   uint16 // PCM = 1
	NumChannels   uint16 // Mono 1, Stereo = 2,
	SampleRate    uint32 // 44100 for CD-Quality
	ByteRate      uint32 // SampleRate * NumChannels * BitsPerSample / 8
	BlockAlign    uint16 // NumChannels * BitsPerSample / 8 (number of bytes per sample)
	BitsPerSample uint16 // 8 bits, 16 bits..
}

type Wav struct {
	header WaveHeader
	fmt    WaveFmt
	data   []byte
}

type WaveHeader struct {
	ChunkID   string
	ChunkSize uint32
	Format    string
}

// Great WAV format description
// http://www-mmsp.ece.mcgill.ca/Documents/AudioFormats/WAVE/WAVE.html
// http://soundfile.sapp.org/doc/WaveFormat/
func LoadWav(path string) (*Wav, error) {
	wav := Wav{}

	file, err := os.Open(path)
	if err != nil {
		return &wav, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return &wav, err
	}

	// Parse Header
	wav.header.ChunkID = string(data[0:4])
	if wav.header.ChunkID != "RIFF" {
		println(fmt.Sprintf("Bad magic number for '%s' (%s)"), path, wav.header.ChunkID)
		os.Exit(1)
	}

	wav.header.ChunkSize = binary.LittleEndian.Uint32(data[4:8])
	if int(wav.header.ChunkSize) != len(data[8:]) {
		println("Bad wav ", path, "size should be", wav.header.ChunkSize, "but is", len(data[10:]))
		os.Exit(1)
	}
	wav.header.Format = string(data[8:12])
	if wav.header.Format != "WAVE" {
		println(fmt.Sprintf("Bad wav format for '%s' (%s)"), path, wav.header.ChunkID)
		os.Exit(1)
	}

	// Now parse "fmt " chunk
	data = data[12:]
	wav.fmt.Subchunk1ID = string(data[0:4])
	if "fmt " != wav.fmt.Subchunk1ID {
		println(fmt.Sprintf("Unexpected subChunk in '%s' (%s)", path, wav.fmt.Subchunk1ID))
		os.Exit(1)
	}

	wav.fmt.Subchunk1Size = binary.LittleEndian.Uint32(data[4:8])
	if 16 != wav.fmt.Subchunk1Size {
		println(fmt.Sprintf("Unexpected Subchunk1Size in '%s' (%d) expected 16", path, wav.fmt.Subchunk1Size))
		os.Exit(1)
	}

	wav.fmt.AudioFormat = binary.LittleEndian.Uint16(data[8:10])
	if 1 != wav.fmt.AudioFormat {
		println(fmt.Sprintf("Unexpected AudioFormat in '%s' (%d)", path, wav.fmt.AudioFormat))
		os.Exit(1)
	}

	wav.fmt.NumChannels = binary.LittleEndian.Uint16(data[10:12])
	if 1 != wav.fmt.NumChannels {
		println(fmt.Sprintf("Unexpected NumChannels in '%s' (%d) MUST be 1", path, wav.fmt.NumChannels))
		os.Exit(1)
	}

	wav.fmt.SampleRate = binary.LittleEndian.Uint32(data[12:16])
	if wav.fmt.SampleRate != 7575 {
		println(fmt.Sprintf("Unexpected sample rate in '%s' (%d) MUST be 7575", path, wav.fmt.SampleRate))
		os.Exit(1)
	}
	wav.fmt.ByteRate = binary.LittleEndian.Uint32(data[16:20])
	wav.fmt.BlockAlign = binary.LittleEndian.Uint16(data[20:22])

	wav.fmt.BitsPerSample = binary.LittleEndian.Uint16(data[22:24])
	if 8 != wav.fmt.BitsPerSample {
		println(fmt.Sprintf("Unexpected NumChannels in '%s' (%d) MUST be 8", path, wav.fmt.BitsPerSample))
		os.Exit(1)
	}

	// Now parse the data chunk
	data = data[24:]
	chunkName := string(data[0:4])
	if chunkName != "data" {
		println(fmt.Sprintf("Unexpexcted chunk '%s' (%s) but expected 'data'"), path, chunkName)
		os.Exit(1)
	}

	dataLength := binary.LittleEndian.Uint32(data[4:8])
	wav.data = data[8:]
	// Not all the rest of the file is PCM. If the payload has uneven bytes number, it
	// is padded with either a 0 or a 1 at the end. Slice accordingly
	if dataLength%2 != 0 {
		wav.data = wav.data[:len(wav.data)-1]
	}
	if dataLength != uint32(len(wav.data)) {
		println(fmt.Sprintf("Bad data chunk size '%d' but expected '%d'"), dataLength, len(wav.data))
		os.Exit(1)
	}

	return &wav, nil
}
