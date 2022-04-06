package sites

import "os"

func ensureDir(dir string) {
	// Test if there is a gfx src folder. If not, return null
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0777)
		if err != nil {
			println("Unable to create folder ", dir, ":", err.Error())
			os.Exit(1)
		}
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
