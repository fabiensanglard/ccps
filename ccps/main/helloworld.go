package main

import (
	"ccps/sites"
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

//go:embed hwAssets/helloworld.png
var helloWorldSprite []byte

//go:embed hwAssets/helloworld.wav
var helloWorldSound []byte

func helloWorld(args []string) {
	post(args)

	fs := flag.NewFlagSet("hwFlags", flag.ContinueOnError)
	v := fs.Bool("v", false, "Verbose mode")
	verbose := *v

	if err := fs.Parse(args); err != nil {
		println(fmt.Sprintf("Cmd parsing error '%s'", err))
		os.Exit(1)
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
