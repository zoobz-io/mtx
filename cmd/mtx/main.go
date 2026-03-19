package main

import (
	"os"

	"github.com/zoobzio/mtx"
)

func main() {
	if err := mtx.Execute(); err != nil {
		os.Exit(1)
	}
}
