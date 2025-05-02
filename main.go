package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/MatusOllah/makesticker/assets"
	cmdcolor "github.com/fatih/color"
	"github.com/fogleman/gg"
	"github.com/mazznoer/csscolorparser"
	"github.com/spf13/pflag"
	"github.com/theckman/yacspin"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

var spin *yacspin.Spinner

func handleError(err error) {
	if err != nil {
		if spin != nil {
			if serr := spin.StopFail(); serr != nil {
				panic(err) // If this fails, blame Tsukasa xD
			}
		}
		red := cmdcolor.New(cmdcolor.FgRed).SprintFunc()
		fmt.Fprintf(os.Stderr, "%s %s\n", red("Error:"), err.Error())
		os.Exit(1)
	}
}

func getCharacterImage(character string) (image.Image, error) {
	// From local disk
	if f, err := os.Open(character); err == nil {
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			return nil, fmt.Errorf("failed to decode character image: %w", err)
		}

		return img, nil
	}

	// From embedded VFS
	f, err := assets.FS.Open("characters/" + character + ".png")
	if err != nil {
		return nil, fmt.Errorf("failed to open character image VFS file: %w", err)
	}

	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode character image: %w", err)
	}

	return img, nil
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
		fmt.Fprintf(os.Stderr, "MakeSticker is a Project SEKAI CLI sticker maker, written in Go.\n")
		fmt.Fprintf(os.Stderr, "Copyright (c) 2025 Matúš Ollah; Licensed under MIT License\n\n")
		pflag.Usage()
		os.Exit(0)
	}

	if *versionFlag {
		pink := cmdcolor.New(cmdcolor.FgHiMagenta, cmdcolor.Bold).SprintFunc()
		cyan := cmdcolor.New(cmdcolor.FgCyan, cmdcolor.Bold).SprintFunc()
		bold := cmdcolor.New(cmdcolor.Bold).SprintFunc()

		if ver := version(); ver != "" {
			fmt.Fprintf(os.Stderr, "%s %s %s\n", pink("MakeSticker"), bold("version"), ver)
		} else {
			fmt.Fprintln(os.Stderr, "No version info available for this build.")
		}
		fmt.Fprintf(os.Stderr, "%s %s %s (%s/%s)\n", cyan("Go"), bold("version"), runtime.Version(), runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	if *listCharsFlag {
		handleError(fs.WalkDir(assets.FS, "characters", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if !d.IsDir() {
				fmt.Fprintln(os.Stderr, strings.ReplaceAll(filepath.Base(path), ".png", ""))
			}

			return nil
		}))
		os.Exit(0)
	}

	// Positional arguments
	if len(pflag.Args()) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] character text\n\n", os.Args[0])
		handleError(fmt.Errorf("invalid arguments"))
	}
	character := pflag.Args()[0]
	text := pflag.Args()[1]

	if !*quietFlag {
		var err error
		spin, err = yacspin.New(yacspin.Config{
			Frequency:         200 * time.Millisecond,
			Writer:            os.Stderr,
			ColorAll:          true,
			CharSet:           yacspin.CharSets[9],
			Message:           " Generating",
			StopFailMessage:   " Failed",
			StopFailCharacter: "✗",
			StopFailColors:    []string{"fgRed"},
		})
		handleError(err)
		handleError(spin.Start())
	}

	// Character image
	characterImg, err := getCharacterImage(character)
	handleError(err)

	// Font
	textFace, err := parseFont(*textFontSizeFlag)
	handleError(err)

	// Text color
	// by the way, I could have just used dc.SetColorHex instead of parsing manually but I wanted more flexibility
	textColor, err := getTextColor(character)
	handleError(err)

	dc := gg.NewContext(296, 256)

	// Character image
	dc.DrawImage(characterImg, (dc.Width()-characterImg.Bounds().Dx())/2, (dc.Height()-characterImg.Bounds().Dy())/2)

	// Text
	dc.SetFontFace(textFace)
	w, h := dc.MeasureString(text)
	dc.Rotate(*textRotationFlag)
	dc.DrawRectangle(*textXFlag, *textYFlag, w, h)
	dc.SetRGB(1, 1, 1)
	drawStringAnchoredOutline(dc, text, *textXFlag, *textYFlag, 0.5, 0.5, 5)
	dc.SetColor(textColor)
	dc.DrawStringAnchored(text, *textXFlag, *textYFlag, 0.5, 0.5)

	// Output
	if *outputFlag == "-" {
		handleError(dc.EncodePNG(os.Stdout))
	} else {
		handleError(dc.SavePNG(*outputFlag))
	}

	if spin != nil {
		handleError(spin.Stop())
	}
}
