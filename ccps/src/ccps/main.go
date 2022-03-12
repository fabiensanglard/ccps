package main

import (
	"fmt"
	"z80"
)
import "boards"

func main()  {
	println("Hello Wordl " + z80.Cc)
	fmt.Printf("%+v", boards.Get("ff"))
}
