package main

import (
	"runtime/debug"
)

var MakeStickerVersion string

func version() string {
	if MakeStickerVersion != "" {
		return MakeStickerVersion
	}

	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	return bi.Main.Version
}
