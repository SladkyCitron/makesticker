package main

import (
	"encoding/json"
	"fmt"
	"image/color"
	"io/fs"
	"os"
	"regexp"
	"runtime"

	"github.com/MatusOllah/makesticker/assets"
	cmdcolor "github.com/fatih/color"
	"github.com/fogleman/gg"
	"github.com/mazznoer/csscolorparser"
	"github.com/spf13/pflag"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

func handleError(err error) {
	if err != nil {
		red := cmdcolor.New(cmdcolor.FgRed, cmdcolor.Bold).SprintFunc()
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

func getTextColor(character string) (color.Color, error) {
	b, err := fs.ReadFile(assets.FS, "text_colors.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open text colors file: %w", err)
	}

	var colors map[string]string
	if err := json.Unmarshal(b, &colors); err != nil {
		return nil, fmt.Errorf("failed to parse text colors JSON: %w", err)
	}

	for pattern, colorString := range colors {
		matched, err := regexp.MatchString(pattern, character)
		if err != nil {
			return nil, fmt.Errorf("failed to match character regex: %w", err)
		}
		if matched {
			c, err := csscolorparser.Parse(colorString)
			if err != nil {
				return nil, fmt.Errorf("failed to parse color string %s: %w", colorString, err)
			}
			return c, nil
		}
	}
	return nil, fmt.Errorf("no matching text color for character: %s", character)
}

func drawStringAnchoredOutline(dc *gg.Context, s string, x float64, y float64, ax float64, ay float64, strokeSize int) {
	for dy := -strokeSize; dy <= strokeSize; dy++ {
		for dx := -strokeSize; dx <= strokeSize; dx++ {
			if dx*dx+dy*dy >= strokeSize*strokeSize {
				// give it rounded corners
				continue
			}
			dc.DrawStringAnchored(s, x+float64(dx), y+float64(dy), ax, ay)
		}
	}
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

	// Positional arguments
	character := pflag.Args()[0] // TODO: the character
	text := pflag.Args()[1]

	// Font
	textFace, err := parseFont(*fontSizeFlag)
	handleError(err)

	// Text color
	// by the way, I could have just used dc.SetColorHex instead of parsing manually but I wanted more flexibility
	textColor, err := getTextColor(character)
	handleError(err)

	dc := gg.NewContext(296, 256)
	dc.SetFontFace(textFace)
	dc.SetRGB(1, 1, 1)
	drawStringAnchoredOutline(dc, text, 128, 128, 0.5, 0.5, 5)
	dc.SetColor(textColor)
	dc.DrawStringAnchored(text, 128, 128, 0.5, 0.5)
	dc.SavePNG(*outputFlag)
}
