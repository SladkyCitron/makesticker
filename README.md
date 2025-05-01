# 🌟 makesticker

**makesticker** is a tiny but mighty CLI tool for creating *Project SEKAI* stickers with your favorite characters and custom text, right from your terminal! 💻🌟

Written in **Go**, robust, and perfect for memes, reactions, or just spreading love in sticker form. 💌

## ✨ Features

* Supports all Project SEKAI characters
* Fully configurable text (position, rotation, size)
* Supports accents aka. diacritics (ľ, š, č, ...)
* Outputs PNG images to file or stdout
* Written in Go and blazingly fast! ⚡

## 📦 Installation

### From binary

See [Releases](https://github.com/MatusOllah/makesticker/releases)

### From source

1. Install Go (<https://go.dev>)
2. Run: `go install -v github.com/MatusOllah/makesticker@latest`

The binary will be installed to `$GOPATH/bin`, so make sure that's in your `$PATH`.

## 🧰 Usage

```sh
makesticker [options] character "your message here"
```

For example:

```sh
makesticker emu_13 "Wonderhoy!" -o my_sticker.png
```

➡️ This produces a sticker with **Emu** and the text "*Wonderhoy!*", saved as `my_sticker.png`.

<p align="center">
    <img src="docs/sample_sticker.png" alt="Wonderhoy!" width="200">
</p>

### Available Characters

You can find the full list of available characters [here](CHARACTERS.md).

Example IDs include: `miku_01`, `emu_13`, `tsukasa_11`, `airi_02`, etc...

### Options

```
  -h, --help                   Print help message
  -o, --output string          Output image (- = stdout) (default "sticker.png")
  -q, --quiet                  Be quiet and disable spinner
      --text-font-size float   Text font size (default 36)
      --text-rotation float    Text rotation in radians (default -0.2)
      --text-x float           Text X position (default 148)
      --text-y float           Text Y position (default 58)
      --version                Print version and exit
```

## 🖼️ Output

Your sticker will be saved as a PNG file, or via stdout, with your chosen character and message.

Perfect for Discord, WhatsApp, Reddit, or anywhere really! 😄

## ⚖️ License

Copyright (c) 2025 Matúš Ollah

Licensed under the **MIT License** (see [LICENSE](LICENSE)) - free to use, fork, remix, and share!
