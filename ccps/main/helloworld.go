package main

import (
	"ccps/sites"
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
)

//go:embed hwAssets/helloworld.png
var helloWorldSprite []byte

//go:embed hwAssets/helloworld.wav
var helloWorldSound []byte

//go:embed hwAssets/m68k/crt0.s
var hwSrcM68kCrt0 []byte

//go:embed hwAssets/m68k/main.c
var hwSrcM68kMain []byte

//go:embed hwAssets/z80/crt0.s
var hwSrcZ80Crt0 []byte

//go:embed hwAssets/z80/main.c
var hwSrcZ80Main []byte

func helloWorld(args []string) {
	postWithBytes(args, hwSrcM68kCrt0, hwSrcM68kMain, hwSrcZ80Crt0, hwSrcZ80Main)

	fs := flag.NewFlagSet("hwFlags", flag.ContinueOnError)
	v := fs.Bool("v", false, "Verbose mode")
	fs.String("b", "", "Target board")
	verbose := *v

	if err := fs.Parse(args); err != nil {
		panic(fmt.Sprintf("Cmd parsing error '%s'", err))
	}

	if verbose {
		println("Starting to generate HelloWorld")
	}

	// Drop a helloWorld GFX
	sites.EnsureDirGFX()
	ioutil.WriteFile(sites.GfxObjPath+"helloworld.png", helloWorldSprite, 0644)

	// Drop a helloWorld SFX
	sites.EnsureDirSFX()
	ioutil.WriteFile(sites.SfxSrcPath+"helloworld.wav", helloWorldSound, 0644)

}
