package locale

import (
	"regexp"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"common/config"
)

var reLocale = regexp.MustCompile(`\[localize\:([a-z\,]*)\:([a-z]+)\]`)

func preprocessingOps(modifiers []string, locale language.Tag) []func(string) string {
	ops := make([]func(string) string, len(modifiers))
	for i, mod := range modifiers {
		switch mod {
		case "capitalize":
			ops[i] = func(s string) string {
				return cases.Title(locale).String(s)
			}
		}
	}
	return ops
}

func PrepareTemplate(localizer config.Localizer, t, locale string) string {
	if locale == "" {
		return t
	}
	for _, match := range reLocale.FindAllStringSubmatch(t, -1) {
		modifiers := match[1]
		entity := localizer.Localize(match[2], locale)
		modOps := preprocessingOps(strings.Split(modifiers, ","), language.Make(locale))
		for _, op := range modOps {
			entity = op(entity)
		}
		t = strings.ReplaceAll(t, match[0], entity)
	}
	return t
}
