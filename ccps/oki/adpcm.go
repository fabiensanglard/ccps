package oki

func CLAMP(x int, low int, high int) int {
	if x > high {
		return high
	}
	if x < low {
		return low
	}
	return x
}

var stepTable = [...]int16{
	16, 17, 19, 21, 23, 25, 28, 31,
	34, 37, 41, 45, 50, 55, 60, 66,
	73, 80, 88, 97, 107, 118, 130, 143,
	157, 173, 190, 209, 230, 253, 279, 307,
	337, 371, 408, 449, 494, 544, 598, 658,
	724, 796, 876, 963, 1060, 1166, 1282, 1411,
	1552}

func oki_step(step byte, lastSample *int16, lastStep *byte) int16 {

	var deltaTable = [...]int16{
		1, 3, 5, 7, 9, 11, 13, 15,
		-1, -3, -5, -7, -9, -11, -13, -15}

	var adjustFactor = [...]int8{
		-1, -1, -1, -1, 2, 4, 6, 8}

	stepSize := stepTable[*lastStep]
	delta := deltaTable[step&15] * stepSize / 8
	out := *lastSample + delta
	out = int16(CLAMP(int(out), -2048, 2047))
	*lastSample = out

	// Adjust step
	adjustedStep := int8(int(*lastStep) + int(adjustFactor[step&7]))
	*lastStep = byte(CLAMP(int(adjustedStep), 0, len(stepTable)-1))

	return out
}

func oki_encode_step(sample byte, lastSample *int16, lastStep *byte) byte {
	step_size := stepTable[*lastStep]
	delta := int16(sample) - *lastSample

	var adpcm_sample byte = 0
	if delta < 0 {
		adpcm_sample = 0x8 // Set bit 4 to 1
	}

	// delta = abs(delta)
	if delta < 0 {
		delta = -delta
	}

	for bit := 2; bit >= 0; bit-- {
		if delta >= step_size {
			adpcm_sample |= (1 << bit)
			delta -= step_size
		}
		step_size >>= 1
	}
	oki_step(adpcm_sample, lastSample, lastStep)
	return adpcm_sample
}

func PCMtoADPCM(wav []byte) []byte {
	lastSample := int16(0)
	lastStep := byte(0)

	if len(wav)%2 == 1 {
		wav = append(wav, wav[len(wav)-1])
	}

	adpcm := make([]byte, len(wav)/2)
	adpcmCursor := 0

	for i := 0; i < len(wav); i += 2 {
		nibble1 := oki_encode_step(wav[i], &lastSample, &lastStep) & 0xF
		nibble2 := oki_encode_step(wav[i+1], &lastSample, &lastStep) & 0xF
		adpcm[adpcmCursor] = nibble1<<4 | nibble2
		adpcmCursor++
	}
	return adpcm
}

func ADPCMToPCM(adpcm []byte) []byte {
	lastSample := int16(0)
	lastStep := byte(0)

	pcm := make([]byte, len(adpcm)*2)
	pcmCursor := 0
	adpcmCursor := 0
	for i := 0; i < len(adpcm); i++ {
		twoNibbles := adpcm[adpcmCursor]
		pcm[pcmCursor] = byte(oki_step(twoNibbles>>4, &lastSample, &lastStep))
		pcmCursor++
		pcm[pcmCursor] = byte(oki_step(twoNibbles&0xF, &lastSample, &lastStep))
		pcmCursor++
	}

	return pcm
}
