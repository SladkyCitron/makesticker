package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	assetsFlag := flag.String("assets", "assets", "Path to assets directory")
	outputFlag := flag.String("output", "CHARACTERS.md", "Output file")
	flag.Parse()

	var mdBuf bytes.Buffer

	fmt.Fprintln(&mdBuf, "# List of characters")
	fmt.Fprintln(&mdBuf)
	fmt.Fprintln(&mdBuf, "`CTRL+F` is your friend! 😉")
	fmt.Fprintln(&mdBuf)
	fmt.Fprintln(&mdBuf, "| Character ID | Sample image |")
	fmt.Fprintln(&mdBuf, "|--------------|--------------|")

	err := filepath.WalkDir(filepath.Join(*assetsFlag, "characters"), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		fmt.Fprintf(&mdBuf, "| %[1]s | <img src=\"assets/characters/%[1]s.png\" alt=\"%[1]s\" width=\"200\"> |\n", strings.ReplaceAll(filepath.Base(path), ".png", ""))

		return nil
	})
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile(*outputFlag, mdBuf.Bytes(), 0o666); err != nil {
		panic(err)
	}
}
