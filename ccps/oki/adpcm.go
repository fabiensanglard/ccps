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

func oki_step(step byte, history *int16, step_hist *byte) int16 {

	var deltaTable = [...]int8{
		1, 3, 5, 7, 9, 11, 13, 15,
		-1, -3, -5, -7, -9, -11, -13, -15}

	var adjustTable = [...]int8{
		-1, -1, -1, -1, 2, 4, 6, 8}

	stepSize := stepTable[*step_hist]
	delta := int16(deltaTable[step&15]) * stepSize / 8
	out := *history + delta
	out = int16(CLAMP(int(out), -2048, 2047))
	*history = out
	adjustedStep := int8(int(*step_hist) + int(adjustTable[step&7]))
	*step_hist = byte(CLAMP(int(adjustedStep), 0, len(stepTable)-1))

	return out
}

func oki_encode_step(sample byte, history *int16, step_hist *byte) byte {
	step_size := stepTable[*step_hist]
	delta := int16(sample) - *history

	var adpcm_sample byte
	if delta < 0 {
		adpcm_sample = 8
	} else {
		adpcm_sample = 0
	}

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
	oki_step(adpcm_sample, history, step_hist)
	return adpcm_sample
}

func PCMtoADPCM(wav []byte) []byte {
	history := int16(0)
	step_hist := byte(0)
	buf_sample := byte(0)
	nibble := 0

	adpcm := make([]byte, len(wav)/2)
	adpcmCursor := 0

	for i := 0; i < len(wav); i++ {
		sample := wav[i]
		step := oki_encode_step(sample, &history, &step_hist)
		if nibble > 0 {
			adpcm[adpcmCursor] = buf_sample | (step & 0xF)
			adpcmCursor++
		} else {
			buf_sample = (step & 0xF) << 4
		}
		nibble ^= 1
	}
	return adpcm
}

func ADPCMToPCM(adpcm []byte) []byte {
	history := int16(0)
	step_hist := byte(0)
	nibble := byte(0)

	pcm := make([]byte, len(adpcm)*2)
	pcmCursor := 0
	adpcmCursor := 0
	for i := 0; i < len(adpcm); i++ {
		step := adpcm[adpcmCursor] << nibble
		step >>= 4
		if nibble != 0 {
			adpcmCursor++
		}
		nibble ^= 4
		pcm[pcmCursor] = byte(oki_step(step, &history, &step_hist))
		pcmCursor++
	}

	return pcm
}
