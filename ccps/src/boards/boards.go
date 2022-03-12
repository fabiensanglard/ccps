package boards

import "fmt"

type CpsBRegisters struct {
   palette int
}

type GFXROM struct {

}

type Board struct{
	cpsBReg CpsBRegisters
	gfx GFXROM
}

func Get(name string) Board {
	if (name == "sf2") {
		return sf2Board()
	}
	panic(fmt.Sprintf("Unknown board '%s'", name))
}

func sf2Board() Board {
  return Board{}
}