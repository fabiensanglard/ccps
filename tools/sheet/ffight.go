package main

type FFight struct {
	Game
}

func (game FFight) GetName() string {
	return game.name
}

func makeFFight() FFight {
	var game FFight
	game.gfxROMSize = 0x80000
	game.gfx_banks = []RomSrc{
		{"ff-5m.7a", "7868f5801347340867720255f8380548ad1a65a7", 2, 0x00000, 0x20000, 0x00000, 8},
		{"ff-7m.9a", "f7b00a3ca8cb85264ab293089f9f540a8292b49c", 2, 0x00000, 0x20000, 0x00002, 8},
		{"ff-1m.3a", "5ce16af72858a57aefbf6efed820c2c51935882a", 2, 0x00000, 0x20000, 0x00004, 8},
		{"ff-3m.5a", "df5f3d3aa96a7a33ff22f2a31382942c4c4f1111", 2, 0x00000, 0x20000, 0x00006, 8},
	}
	game.name = "ffight"
	game.paletteAddr = 0

	return game
}
