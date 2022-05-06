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
	data   []int16
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
		panic(fmt.Sprintf("Bad magic number for '%s' (%s)", path, wav.header.ChunkID))
	}

	wav.header.ChunkSize = binary.LittleEndian.Uint32(data[4:8])
	if int(wav.header.ChunkSize) != len(data[8:]) {
		panic(fmt.Sprintf("Bad wav '%s' size should be %d but is %d", path, wav.header.ChunkSize, len(data[10:])))

	}
	wav.header.Format = string(data[8:12])
	if wav.header.Format != "WAVE" {
		panic(fmt.Sprintf("Bad wav format for '%s' (%s)", path, wav.header.ChunkID))
	}

	// Now parse "fmt " chunk
	data = data[12:]
	wav.fmt.Subchunk1ID = string(data[0:4])
	if "fmt " != wav.fmt.Subchunk1ID {
		panic(fmt.Sprintf("Unexpected subChunk in '%s' (%s)", path, wav.fmt.Subchunk1ID))
	}

	wav.fmt.Subchunk1Size = binary.LittleEndian.Uint32(data[4:8])
	if 16 != wav.fmt.Subchunk1Size {
		panic(fmt.Sprintf("Unexpected Subchunk1Size in '%s' (%d) expected 16", path, wav.fmt.Subchunk1Size))
	}

	wav.fmt.AudioFormat = binary.LittleEndian.Uint16(data[8:10])
	if 1 != wav.fmt.AudioFormat {
		panic(fmt.Sprintf("Unexpected AudioFormat in '%s' (%d)", path, wav.fmt.AudioFormat))
	}

	wav.fmt.NumChannels = binary.LittleEndian.Uint16(data[10:12])
	if 1 != wav.fmt.NumChannels {
		panic(fmt.Sprintf("Unexpected NumChannels in '%s' (%d) MUST be 1", path, wav.fmt.NumChannels))
	}

	wav.fmt.SampleRate = binary.LittleEndian.Uint32(data[12:16])
	if wav.fmt.SampleRate != 7575 {
		panic(fmt.Sprintf("Unexpected sample rate in '%s' (%d) MUST be 7575", path, wav.fmt.SampleRate))
	}
	wav.fmt.ByteRate = binary.LittleEndian.Uint32(data[16:20])
	wav.fmt.BlockAlign = binary.LittleEndian.Uint16(data[20:22])

	wav.fmt.BitsPerSample = binary.LittleEndian.Uint16(data[22:24])
	if 16 != wav.fmt.BitsPerSample {
		panic(fmt.Sprintf("Unexpected BitPerSample in '%s' (%d) MUST be 16", path, wav.fmt.BitsPerSample))
	}

	// Now parse the data chunk
	data = data[24:]
	chunkName := string(data[0:4])
	if chunkName != "data" {
		panic(fmt.Sprintf("Unexpexcted chunk '%s' (%s) but expected 'data'", path, chunkName))
	}

	dataLength := binary.LittleEndian.Uint32(data[4:8])
	payload := data[8:]
	if len(payload)%2 == 1 {
		payload = append(payload, payload[len(payload)-1])
	}

	wav.data = toArray16(payload)
	// Not all the rest of the file is PCM. If the payload has uneven bytes number, it
	// is padded with either a 0 or a 1 at the end. Slice accordingly
	if dataLength%2 != 0 {
		wav.data = wav.data[:len(wav.data)-1]
	}
	//if dataLength != uint32(len(wav.data)) {
	//	println(fmt.Sprintf("Bad data chunk size '%d' but expected '%d'"), dataLength, len(wav.data))
	//	os.Exit(1)
	//}

	return &wav, nil
}

func toArray16(payload []byte) []int16 {
	wav := make([]int16, len(payload)/2)
	cursor := 0
	for i := 0; i < len(payload)-3; i += 2 {
		wav[cursor] = int16(binary.LittleEndian.Uint16(payload[i : i+2]))
		cursor++
	}
	return wav
}
