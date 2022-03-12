package ccz80

import (
	"fmt"
	"os"
)

/*
   .cps/config.txt
   .cps/z80
   .cps/68000
   out
   code/z80
   code/68000
   gfx/*.png
   sfx/*.wav
   sfx/*.vgm
*/
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("ccps [CMD] args ...")
	fmt.Println("  init")
	fmt.Println("  cc [CPU]")
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		printUsage()
		return
	}
	fmt.Println("Hello World")
}
