package main

import (
	"fmt"
	"runtime/debug"
	"time"
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

	// If no main version is available, Go defaults to (devel)
	if bi.Main.Version != "(devel)" {
		return bi.Main.Version
	}

	var vcsRevision string
	var vcsTime time.Time
	for _, setting := range bi.Settings {
		switch setting.Key {
		case "vcs.revision":
			vcsRevision = setting.Value
		case "vcs.time":
			vcsTime, _ = time.Parse(time.RFC3339, setting.Value)
		}
	}

	if vcsRevision != "" {
		return fmt.Sprintf("%s, (%s)", vcsRevision, vcsTime)
	}

	return ""
}
