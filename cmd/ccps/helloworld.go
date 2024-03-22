package main

import (
	_ "embed"
	"os"

	"github.com/fabiensanglard/ccps/sites"
	"github.com/spf13/cobra"
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

func helloWorld(cmd *cobra.Command, args []string) {
	postWithBytes(cmd, hwSrcM68kCrt0, hwSrcM68kMain, hwSrcZ80Crt0, hwSrcZ80Main)

	if verbose {
		cmd.Println("Starting to generate HelloWorld")
	}

	// Drop a helloWorld GFX
	sites.EnsureDirGFX()
	if err := os.WriteFile(sites.GfxObjPath+"helloworld.png", helloWorldSprite, 0644); err != nil {
		cmd.PrintErr(err)
		os.Exit(1)
	}

	// Drop a helloWorld SFX
	sites.EnsureDirSFX()
	if err := os.WriteFile(sites.SfxSrcPath+"helloworld.wav", helloWorldSound, 0644); err != nil {
		cmd.PrintErr(err)
		os.Exit(1)
	}

}
