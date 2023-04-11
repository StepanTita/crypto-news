package config

import (
	"io/fs"
	"strings"

	"github.com/pkg/errors"

	"templates"
)

type Templator interface {
	Template(name string) string
}

type templator struct {
	templates map[string]string
}

func NewTemplator() Templator {
	templatesMapping := make(map[string]string)

	err := fs.WalkDir(templates.Dir, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err, "failed to walk templates dir")
		}

		if d.IsDir() {
			return nil
		}

		// expect name.command.tmpl
		commandName := strings.Split(d.Name(), ".")[0]
		rawContent, err := fs.ReadFile(templates.Dir, d.Name())
		if err != nil {
			return errors.Wrap(err, "failed to read file")
		}
		templatesMapping[commandName] = string(rawContent)

		return nil
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to walk dir"))
	}

	return &templator{
		templates: templatesMapping,
	}
}

func (l templator) Template(name string) string {
	return l.templates[name]
}
