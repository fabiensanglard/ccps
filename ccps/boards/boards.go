package boards

import (
	"fmt"
	"io/ioutil"
	"os"
)

type Rom2Eprom func(string, string)

type CpsBRegisters struct {
	palette int
}

type ROMSet struct {
	Size    int64
	Epromer Rom2Eprom
}

type Board struct {
	CpsBReg CpsBRegisters
	GFX     ROMSet
	Z80     ROMSet
	M68k    ROMSet
	Oki     ROMSet
}

func Get(name string) *Board {
	if name == "sf2" {
		return sf2Board()
	}
	println(fmt.Sprintln("Unknown board '%s'", name))
	os.Exit(1)
	return &Board{}
}

var boards []Board

func sf2Board() *Board {
	sf2 := Board{}
	sf2.Z80.Size = 65536
	sf2.Z80.Epromer = sf2Z80EPromer

	sf2.M68k.Size = 1048576
	sf2.M68k.Epromer = sf2M68kEPromer

	sf2.GFX.Size = 6291456
	sf2.GFX.Epromer = sf2GFXEpromer
	return &sf2
}

func cp(src string, dst string) {
	//cmd := fmt.Sprintf("cp %s %s", src, dst)
	//if verbose
	//println(cmd)

	input, err := ioutil.ReadFile(src)
	if err != nil {
		fmt.Println(fmt.Sprintf("cp error: Cannot open src: '%s'", src))
		os.Exit(1)
	}

	err = ioutil.WriteFile(dst, input, 0644)
	if err != nil {
		fmt.Println(fmt.Sprintf("cp error: Cannot dst: '%s'", dst))
		os.Exit(1)
	}
}

func sf2GFXEpromer(romPath string, outDir string) {
	_, err := ioutil.ReadFile(romPath)

	if err != nil {
		fmt.Println(fmt.Sprintf("Cannot open GFX ROM '%s'", romPath))
		os.Exit(1)
	}

	// Split it
	const BANK_SIZE = 0x200000
	//bank0 := gfxrom[0*BANK_SIZE : 0*BANK_SIZE+BANK_SIZE]
	//bank1 := gfxrom[1*BANK_SIZE : 1*BANK_SIZE+BANK_SIZE]
	//bank2 := gfxrom[2*BANK_SIZE : 2*BANK_SIZE+BANK_SIZE]

	const ROM_SIZE = 0x80000

	//writeToFile(bank0[0:0x80000], 2, 8, ROM_SIZE, outDir+"sf2_06.bin")
	//writeToFile(bank0[2:0x80000], 2, 8, ROM_SIZE, outDir+"sf2_08.bin")
	//writeToFile(bank0[4:0x80000], 2, 8, ROM_SIZE, outDir+"sf2_05.bin")
	//writeToFile(bank0[6:0x80000], 2, 8, ROM_SIZE, outDir+"sf2_07.bin")

	//{"sf2_06.bin", "b9194fb337b30502c1c9501cd6c64ae4035544d4", 2, 0, 0x80000, 0x0000000, 8},
	//{"sf2_08.bin", "3759b851ac0904ec79cbb67a2264d384b6f2f9f9", 2, 0, 0x80000, 0x0000002, 8},
	//{"sf2_05.bin", "520840d727161cf09ca784919fa37bc9b54cc3ce", 2, 0, 0x80000, 0x0000004, 8},
	//{"sf2_07.bin", "2360cff890551f76775739e2d6563858bff80e41", 2, 0, 0x80000, 0x0000006, 8},

	//{"sf2_15.bin", "357c2275af9133fd0bd6fbb1fa9ad5e0b490b3a2", 2, 0, 0x80000, 0x200000, 8},
	//{"sf2_17.bin", "baa92b91cf616bc9e2a8a66adc777ffbf962a51b", 2, 0, 0x80000, 0x200002, 8},
	//{"sf2_14.bin", "2eea16673e60ba7a10bd4d8f6c217bb2441a5b0e", 2, 0, 0x80000, 0x200004, 8},
	//{"sf2_16.bin", "f787aab98668d4c2c54fc4ba677c0cb808e4f31e", 2, 0, 0x80000, 0x200006, 8},

	//{"sf2_25.bin", "5669b845f624b10e7be56bfc89b76592258ce48b", 2, 0, 0x80000, 0x400000, 8},
	//{"sf2_27.bin", "9af9df0826988872662753e9717c48d46f2974b0", 2, 0, 0x80000, 0x400002, 8},
	//{"sf2_24.bin", "a6a7f4725e52678cbd8d557285c01cdccb2c2602", 2, 0, 0x80000, 0x400004, 8},
	//{"sf2_26.bin", "f9a92d614e8877d648449de2612fc8b43c85e4c2", 2, 0, 0x80000, 0x400006, 8},
}

func sf2Z80EPromer(rom string, outputDir string) {
	cp(rom, outputDir+"sf2_9.12a")
}

func sf2M68kEPromer(rom string, outputDir string) {
	r, err := ioutil.ReadFile(rom)

	if err != nil {
		fmt.Println(fmt.Sprintf("Cannot open Z80 ROM '%s'", rom))
		os.Exit(1)
	}

	const ROM_SIZE = 131072
	writeToFile(r[0*ROM_SIZE:], 1, 2, ROM_SIZE, outputDir+"sf2e_30g.11e")
	writeToFile(r[0*ROM_SIZE+1:], 1, 2, ROM_SIZE, outputDir+"sf2e_37g.11f")
	writeToFile(r[1*ROM_SIZE:], 1, 2, ROM_SIZE, outputDir+"sf2e_31g.12e")
	writeToFile(r[1*ROM_SIZE+1:], 1, 2, ROM_SIZE, outputDir+"sf2e_38g.12f")
	writeToFile(r[2*ROM_SIZE:], 1, 2, ROM_SIZE, outputDir+"sf2e_28g.9e")
	writeToFile(r[2*ROM_SIZE+1:], 1, 2, ROM_SIZE, outputDir+"sf2e_35g.9f")
	writeToFile(r[3*ROM_SIZE:], 1, 2, ROM_SIZE, outputDir+"sf2_29b.10e")
	writeToFile(r[3*ROM_SIZE+1:], 1, 2, ROM_SIZE, outputDir+"sf2_36b.10f")
}

func writeToFile(rom []byte, size int, skip int, num int, filename string) {
	var eprom = make([]byte, num)

	for i := 0; i < num; i++ {
		eprom[i] = rom[i*skip]
	}
	err := ioutil.WriteFile(filename, eprom, 0644)
	if err != nil {
		fmt.Println("Unable to write EPROM '", filename, "'")
		os.Exit(1)
	}
}
