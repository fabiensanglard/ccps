package main

import (
	"ccps/boards"
	"ccps/z80"
	"fmt"
)

func main() {
	println("Starting...")

	z80Roms := z80.Build()
	for _, rom := range z80Roms {
		println(rom)
	}
	fmt.Printf("%+v", boards.Get("ff"))
}
