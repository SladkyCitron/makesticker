package main

import (
	"fmt"
	"io/fs"
	"os"
	"runtime"

	"github.com/MatusOllah/makesticker/assets"
	"github.com/fatih/color"
	"github.com/fogleman/gg"
	"github.com/spf13/pflag"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func handleError(err error) {
	if err != nil {
		red := color.New(color.FgRed, color.Bold).SprintFunc()
		fmt.Fprintf(os.Stderr, "%s %s\n", red("Error:"), err.Error())
		os.Exit(1)
	}
}

func parseFont(size float64) (font.Face, error) {
	b, err := fs.ReadFile(assets.FS, "FOT-YurukaStd.otf")
	if err != nil {
		return nil, fmt.Errorf("failed to open font file: %w", err)
	}

	font, err := opentype.Parse(b)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	return opentype.NewFace(font, &opentype.FaceOptions{Size: size, DPI: 72})
}

func main() {
	if *helpFlag {
		fmt.Fprintf(os.Stderr, "makesticker is a Project SEKAI CLI sticker maker, written in Go.\n\n")
		pflag.Usage()
		os.Exit(0)
	}

	if *versionFlag {
		fmt.Fprintf(os.Stderr, "makesticker version %s\n", Version)
		fmt.Fprintf(os.Stderr, "Go version %s (%s/%s)\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	if len(pflag.Args()) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] character text\n\n", os.Args[0])
		handleError(fmt.Errorf("invalid arguments"))
	}

	_ = pflag.Args()[0] // TODO: the character
	text := pflag.Args()[1]

	textFace, err := parseFont(*fontSizeFlag)
	handleError(err)

	ctx := gg.NewContext(296, 256)
	ctx.SetFontFace(textFace)
	ctx.SetRGB(0, 0, 0)
	ctx.DrawStringAnchored(text, 128, 128, 0.5, 0.5)
	ctx.SavePNG(*outputFlag)
}
