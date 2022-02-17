package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

const ROM_SIZE = 131072

func makeRom(in []byte, name string) {
	rom := make([]byte, ROM_SIZE)

	for i := 0; i < ROM_SIZE; i += 1 {
		rom[i] = in[i*2]
	}

	err := ioutil.WriteFile(name, rom, 0644)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func outputSet(in []byte, o1 string, o2 string) {
	makeRom(in[0:], o1)
	makeRom(in[1:], o2)
}

func main() {
	file, err := ioutil.ReadFile("game.rom")

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	outputSet(file[ROM_SIZE*2*0:ROM_SIZE*2*0+ROM_SIZE*2], "sf2e_30g.11e", "sf2e_37g.11f")
	outputSet(file[ROM_SIZE*2*1:ROM_SIZE*2*1+ROM_SIZE*2], "sf2e_31g.12e", "sf2e_38g.12f")
	outputSet(file[ROM_SIZE*2*2:ROM_SIZE*2*2+ROM_SIZE*2], "sf2e_28g.9e", "sf2e_35g.9f")
	outputSet(file[ROM_SIZE*2*3:ROM_SIZE*2*3+ROM_SIZE*2], "sf2_29b.10e", "sf2_36b.10f")

}
