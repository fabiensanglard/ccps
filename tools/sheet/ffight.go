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
		{"ff-5m.7a", "7868f5801347340867720255f8380548ad1a65a7", 2, 0x00000, 0x80000, 0x00000, 8},
		{"ff-7m.9a", "f7b00a3ca8cb85264ab293089f9f540a8292b49c", 2, 0x00000, 0x80000, 0x00002, 8},
		{"ff-1m.3a", "5ce16af72858a57aefbf6efed820c2c51935882a", 2, 0x00000, 0x80000, 0x00004, 8},
		{"ff-3m.5a", "df5f3d3aa96a7a33ff22f2a31382942c4c4f1111", 2, 0x00000, 0x80000, 0x00006, 8},
	}
	game.name = "ffight"

	game.codeROMSize = 0x400000
	game.code_banks = []RomSrc{
		{"ff_36.11f", "0756ae576a1f6d5b8b22f8630dca40ef38567ea6", 1, 0, 0x20000, 0x00000, 2},
		{"ff_42.11h", "5045a467f3e228c02b4a355b52f58263ffa90113", 1, 0, 0x20000, 0x00001, 2},

		{"ff_37.12f", "38f44434c8befd623953ae23d6e5ff4e201d6627", 1, 0, 0x20000, 0x40000, 2},
		{"ffe_43.12h", "de16873d1639ac1738be0937270b108a9914f263", 1, 0, 0x20000, 0x40001, 2},

		{"ff-32m.8h", "d3362dadded31ccb7eaf71ef282d698d18edd722", 1, 0, 0x80000, 0x80000, 1},
	}
	game.paletteAddr = 0x8ACBA

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

func (game *FFight) Load() bool {
	if !game.Game.Load() {
		return false
	}

	// Swap the latest part 16-bit WORDS.
	for i := 0; i < 0x80000-1; i++ {
		v1 := game.codeROM[0x80000+i+0]
		v2 := game.codeROM[0x80000+1+1]
		game.codeROM[0x80000+1+1] = v2
		game.codeROM[0x80000+1+0] = v1

	}

	return true
}
