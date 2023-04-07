package config

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

type Templator interface {
	Template(name string) string
}

type templator struct {
	templates map[string]string
}

func NewTemplator(templatesDir string) Templator {
	templates := make(map[string]string)

	err := filepath.WalkDir(templatesDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return errors.Wrap(err, "failed to walk templates dir")
		}

		if d.IsDir() {
			return nil
		}

		// expect name.command.tmpl
		commandName := strings.Split(d.Name(), ".")[0]
		rawContent, err := os.ReadFile(path)
		if err != nil {
			return errors.Wrap(err, "failed to read file")
		}
		templates[commandName] = string(rawContent)

		return nil
	})
	if err != nil {
		panic(errors.Wrap(err, "failed to walk dir"))
	}

	return &templator{
		templates: templates,
	}
}

func (l templator) Template(name string) string {
	return l.templates[name]
}
