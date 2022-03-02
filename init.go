package main

import "os"

func init() {
	err := os.MkdirAll(CONF_DIR, os.ModePerm)
}
