package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

var (
	helpFlag         = pflag.BoolP("help", "h", false, "Print help message")
	versionFlag      = pflag.Bool("version", false, "Print version and exit")
	outputFlag       = pflag.StringP("output", "o", "sticker.png", "Output image (- = stdout)")
	textFontSizeFlag = pflag.Float64("text-font-size", 36, "Text font size")
	textXFlag        = pflag.Float64("text-x", 148, "Text X position")
	textYFlag        = pflag.Float64("text-y", 58, "Text Y position")
	textRotationFlag = pflag.Float64("text-rotation", -0.2, "Text rotation in radians")
)

func init() {
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] character text\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		pflag.PrintDefaults()
	}

	pflag.Parse()
}
