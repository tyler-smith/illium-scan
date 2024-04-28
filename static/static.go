package static

import "embed"

//go:embed all:css
var StaticFiles embed.FS
