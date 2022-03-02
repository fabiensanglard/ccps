package main

const BASE string = "./"
const CONF_DIR string = BASE + ".cps/"
const CONF_FILE string = CONF_DIR + "config.txt"

type Config struct {
	cc_z80 string
	cc_68k string
	cpsb   int
}

func Config_Load() Config {
	var conf Config
	conf.cpsb = 11
	return conf
}
