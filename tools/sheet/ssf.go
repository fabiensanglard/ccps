package main

type SSF struct {
	CPS2Game
}

func (game *SSF) GetName() string {
	return game.name
}

func makeSSF() SSF {
	var game SSF
	game.gfxROMSize = 0xc00000
	game.gfx_banks = []RomSrc{
		{"ssf.13m", "bf2a6d98a656d1cb5734da7836686242d3211137", 2, 0x000000, 0x200000, 0x000000, 8},
		{"ssf.15m", "4b302dbb66e8a5c2ad92798699391e981bada427", 2, 0x000000, 0x200000, 0x000002, 8},
		{"ssf.17m", "b21b1c749a8241440879bf8e7cb33968ccef97e5", 2, 0x000000, 0x200000, 0x000004, 8},
		{"ssf.19m", "4d320fc96d1ef0b9928a8ce801734245a4c097a5", 2, 0x000000, 0x200000, 0x000006, 8},
		{"ssf.14m", "0f4d26af338dab5dce5b7b34d32ad0c573434ace", 2, 0x000000, 0x100000, 0x800000, 8},
		{"ssf.16m", "f4456833fb396e6501f4174c0fe5fd63ea40a188", 2, 0x000000, 0x100000, 0x800002, 8},
		{"ssf.18m", "4b060501e56b9d61294748da5387cdae5280ec4d", 2, 0x000000, 0x100000, 0x800004, 8},
		{"ssf.20m", "32b11ba7a12004aff810d719bff7508204c7b7c0", 2, 0x000000, 0x100000, 0x800006, 8},
	}

	game.paletteAddr = 0

	game.name = "ssf"
	return game
}
