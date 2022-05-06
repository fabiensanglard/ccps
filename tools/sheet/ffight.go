package main

type FFight struct {
	Game
}

func (game FFight) GetName() string {
	return game.name
}

func makeFFight() FFight {
	var game FFight
	game.gfxROMSize = 0x200000
	game.gfx_banks = []RomSrc{
		{"ff-5m.7a", "7868f5801347340867720255f8380548ad1a65a7", 2, 0x00000, 0x80000 , 0x00000, 8},
		{"ff-7m.9a", "f7b00a3ca8cb85264ab293089f9f540a8292b49c", 2, 0x00000, 0x80000 , 0x00002, 8},
		{"ff-1m.3a", "5ce16af72858a57aefbf6efed820c2c51935882a", 2, 0x00000, 0x80000 , 0x00004, 8},
		{"ff-3m.5a", "df5f3d3aa96a7a33ff22f2a31382942c4c4f1111", 2, 0x00000, 0x80000 , 0x00006, 8},
	}
	game.name = "ffight"
	game.paletteAddr = 0

	game.codeROMSize = 0x100000
	game.code_banks = []RomSrc{
		{"sf2e_30g.11e", "22558eb15e035b09b80935a32b8425d91cd79669", 1, 0, 0x20000, 0x00000, 2},
		{"sf2e_37g.11f", "bf1ccfe7cc1133f0f65556430311108722add1f2", 1, 0, 0x20000, 0x00001, 2},

		{"sf2e_31g.12e", "86a3954335310865b14ce8b4e0e4499feb14fc12", 1, 0, 0x20000, 0x40000, 2},
		{"sf2e_38g.12f", "6565946591a18eaf46f04c1aa449ee0ae9ac2901", 1, 0, 0x20000, 0x40001, 2},

		{"sf2e_28g.9e", "bbcef63f35e5bff3f373968ba1278dd6bd86b593", 1, 0, 0x20000, 0x80000, 2},
		{"sf2e_35g.9f", "507bda3e4519de237aca919cf72e543403ec9724", 1, 0, 0x20000, 0x80001, 2},

		{"sf2_29b.10e", "75f0827f4f7e9f292add46467f8d4fe19b2514c9", 1, 0, 0x20000, 0xc0000, 2},
		{"sf2_36b.10f", "b807cc495bff3f95d03b061fc629c95f965cb6d8", 1, 0, 0x20000, 0xc0000, 2},
	}
	game.paletteAddr = 0x8ACBA

	return game
}
