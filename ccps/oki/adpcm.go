package oki

var stepSizes = [...]int16{
	16, 17, 19, 21, 23, 25, 28, 31, 34, 37,
	41, 45, 50, 55, 60, 66, 73, 80, 88, 97,
	107, 118, 130, 143, 157, 173, 190, 209, 230, 253,
	279, 307, 337, 371, 408, 449, 494, 544, 598, 658,
	724, 796, 876, 963, 1060, 1166, 1282, 1411, 1552}

var adjustFactor = [...]int16{
	-1, -1, -1, -1, 2, 4, 6, 8}

type codec struct {
	lastSample int16
	stepIndex  int16
}

func (c *codec) encodeStep(sample int16) byte {
	ss := stepSizes[c.stepIndex]
	diff := sample - c.lastSample
	code := byte(0x0)
	sample >>= 3 // Algo uses 12-bit input

	if diff < 0 {
		diff = -diff
		code = 0x8
	}

	if diff >= ss {
		diff -= ss
		code |= 0x4
	}

	if diff >= (ss >> 1) {
		diff -= ss >> 1
		code |= 0x2
	}

	if diff >= (ss >> 2) {
		code |= 0x1
	}

	c.lastSample = c.decodeStep(code)

	return code
}

func (c *codec) encode(pcm []int16) []byte {
	cursor := 0
	adpcm := make([]byte, len(pcm)/2)
	for i := 0; i < len(pcm); i += 2 {
		msb := c.encodeStep(pcm[i]) & 0xF
		lsb := c.encodeStep(pcm[i+1]) & 0xF
		adpcm[cursor] = (msb << 4) | lsb
		cursor++
	}
	return adpcm
}

func (c *codec) decodeStep(code byte) int16 {
	ss := stepSizes[c.stepIndex]
	delta := ((int16(code&0x7)*2 + 1) * ss) >> 3

	if code&0x8 != 0 {
		delta = -delta
	}

	sample := int16(c.lastSample) + delta

	if sample < -2048 {
		sample = -2046
	}
	if sample > 2047 {
		sample = 2047
	}

	c.nextStepIdx(code)
	c.lastSample = sample
	return sample
}

func (c *codec) nextStepIdx(code byte) {
	c.stepIndex += adjustFactor[code&0x7]

	if c.stepIndex < 0 {
		c.stepIndex = 0
	}

	if c.stepIndex > 48 {
		c.stepIndex = 48
	}
}

func (c *codec) decode(adpcm []byte) []int16 {
	pcm := make([]int16, len(adpcm)*2)
	cursor := 0
	for i := 0; i < len(adpcm); i++ {
		pcm[cursor] = c.decodeStep(adpcm[i] >> 4)
		cursor += 1
		pcm[cursor] = c.decodeStep(adpcm[i] & 0xf)
		cursor += 1
	}
	return pcm
}

// https://github.com/nth-eye/vox/blob/main/src/vox.c

func PCMtoADPCM(wav []int16) []byte {
	var adpcm = codec{}
	return adpcm.encode(wav)
}

func ADPCMToPCM(adpcm []byte) []int16 {
	var codec = codec{}
	return codec.decode(adpcm)
}
