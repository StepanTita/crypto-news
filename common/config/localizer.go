package config

import (
	"io/fs"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/crypto-news/localization"
)

type Localizer interface {
	Localize(word, locale string) string
}

type localizer struct {
	// locale -> word -> translation
	mapping map[string]map[string]string
}

func NewLocalizer() Localizer {
	localizationMapping := make(map[string]map[string]string)

	err := fs.WalkDir(localization.Dir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err, "failed to walk templates dir")
		}

		if d.IsDir() {
			return nil
		}

		rawContent, err := fs.ReadFile(localization.Dir, d.Name())
		if err != nil {
			return errors.Wrap(err, "failed to read file")
		}

		localeTag := strings.Split(d.Name(), ".")[0]
		var body map[string]string
		if err := yaml.Unmarshal(rawContent, &body); err != nil {
			return errors.Wrap(err, "failed to unmarshal localization file")
		}
		localizationMapping[localeTag] = body

		return nil
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to walk dir"))
	}

	return &localizer{
		mapping: localizationMapping,
	}
}

func (l localizer) Localize(word, locale string) string {
	return l.mapping[locale][word]
}
