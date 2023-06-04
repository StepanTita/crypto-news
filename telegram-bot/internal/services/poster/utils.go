package poster

import (
	"strings"
)

func escapeKeepingHTML(text string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
	)

	return replacer.Replace(text)
}
