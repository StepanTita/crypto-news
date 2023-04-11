package templates

import "embed"

//go:embed *.tmpl
var Dir embed.FS
