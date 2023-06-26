package transform

import (
	"regexp"
	"strings"

	"golang.org/x/exp/slices"

	"common/convert"
)

var cleanHtmlRegex = regexp.MustCompile(`<.*?>`)
var allHtmlTagsRegex = regexp.MustCompile(`<([A-Za-z][A-Za-z0-9]*)[^>]*>(.*?)</[A-Za-z][A-Za-z0-9]*>`)

// StripHtmlRegex This method uses a regular expresion to remove HTML tags.
func StripHtmlRegex(s *string) *string {
	if s == nil {
		return nil
	}
	return convert.ToPtr(cleanHtmlRegex.ReplaceAllString(*s, ""))
}

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

	return output
}
