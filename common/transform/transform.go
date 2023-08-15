package transform

import (
	"regexp"
	"strings"

	"golang.org/x/exp/slices"
)

var allHtmlTagsRegex = regexp.MustCompile(`<([A-Za-z][A-Za-z0-9]*)[^>]*>(.*?)</[A-Za-z][A-Za-z0-9]*>`)

func CleanUnsupportedHTML(input string) string {
	matches := allHtmlTagsRegex.FindAllStringSubmatch(input, -1)

	output := input
	for _, tag := range matches {
		tagName := tag[1]
		if slices.Contains([]string{"b", "i", "a", "code", "pre"}, tagName) {
			continue
		}
		output = strings.ReplaceAll(output, tag[0], "")
	}

	return strings.ReplaceAll(output, "<br>", "\n")
}
