package sites

import (
	"fmt"
	"os"
)

func ensureDir(dir string) {
	err := os.RemoveAll(dir)
	err = os.MkdirAll(dir, 0777)
	if err != nil {
		panic(fmt.Sprintf("Unable to create dir '%s'", dir))
	}
}

const GfxSrcPath = "gfx/"
const GfxObjPath = GfxSrcPath + "obj/"
const GfxSC1Path = GfxSrcPath + "scr1/"
const GfxSC2Path = GfxSrcPath + "scr2/"
const GfxSc3Path = GfxSrcPath + "scr3/"

func EnsureDirGFX() {
	ensureDir(GfxSrcPath)
	for _, path := range GfxLayersPath {
		ensureDir(path)
	}
}

var GfxLayersPath = [4]string{
	GfxObjPath,
	GfxSC1Path,
	GfxSC2Path,
	GfxSc3Path,
}

const SfxSrcPath = "sfx/"

func EnsureDirSFX() {
	ensureDir(SfxSrcPath)
}

const CodeGenDir = "gen/"
const M68kGenDir = CodeGenDir + "m68k/"
const Z80GenDir = CodeGenDir + "z80/"

func EnsureCodeGenDirs() {
	ensureDir(M68kGenDir)
	ensureDir(Z80GenDir)
}

// OutDir The directory where finalized ROM file are stored.
const OutDir = "out/"

func EnsureOutDir() {
	ensureDir(OutDir)
}

const srcsDir = "src/"
const tmpDir = "objs/"
const M68kSrcsDir = srcsDir + "m68k/"
const M68kObjsDir = tmpDir + "m68k/"
const Z80SrcsDir = srcsDir + "z80/"
const Z80ObjsDir = tmpDir + "z80/"

func EnsureZ80ObjsDir() {
	ensureDir(Z80ObjsDir)
}

func EnsureM68kObjsDir() {
	ensureDir(M68kObjsDir)
}
