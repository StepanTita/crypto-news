package localization

import "embed"

//go:embed *.locale.yaml
var Dir embed.FS
