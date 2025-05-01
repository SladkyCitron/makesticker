package assets

import "embed"

//go:embed [^_]*
var FS embed.FS
