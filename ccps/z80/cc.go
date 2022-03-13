package z80

import "fmt"

const cc = "sdcc"
const as = "sdasz80"
const linker = "sdldz80"
const crt0 = "crt0.s"

func Build() []string {
	var romsNames []string
	romsNames = append(romsNames, "myrom.rom")
	return romsNames
}

func assemble(srcs string) (error, string) {
	cmd := fmt.Sprintf("%s -plogff -o crt0.rel %s",
		as,
		crt0)
	fmt.Printf(cmd)

	return nil, ""
}

func compile(srcs []string) (error, []string) {
	var outputs []string

	for _, src := range srcs {
		output := ""
		cmd := fmt.Sprintf("%s --compile-only -mz80 --data-loc 0xd000 --no-std-crt0 -o %s %s",
			cc,
			output,
			src)
		fmt.Printf(cmd)
		outputs = append(outputs, "")
	}
	return nil, outputs
}

func link() (error, string) {
	return nil, "output"
}
