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

func EnsureDirGFX() {
	ensureDir(GfxSrcPath)
}

const SfxSrcPath = "sfx/"

func EnsureDirSFX() {
	ensureDir(SfxSrcPath)
}
