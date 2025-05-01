package main

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

var (
	helpFlag     = pflag.BoolP("help", "h", false, "Print help message")
	versionFlag  = pflag.Bool("version", false, "Print version and exit")
	outputFlag   = pflag.StringP("output", "o", "sticker.png", "Output image (- = stdout)")
	fontSizeFlag = pflag.Float64("font-size", 12, "Text font size")
)

func init() {
	pflag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] character text\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		pflag.PrintDefaults()
	}

	pflag.Parse()
}
