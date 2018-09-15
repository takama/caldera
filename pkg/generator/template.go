package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/takama/caldera/pkg/config"
	"github.com/takama/caldera/pkg/helper"
)

// render executes templates in service directory with configured data
func render(cfg *config.Config) error {
	return filepath.Walk(
		cfg.Directories.Service,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("Failed to scan service directory: %s", err)
			}
			if info.IsDir() {
				return nil
			}
			name := filepath.Base(path)
			tpl, err := template.New(name).
				Delims("{{[", "]}}").
				Funcs(
					template.FuncMap{
						"toUpper": strings.ToUpper,
					},
				).ParseFiles(path)
			if err != nil {
				return fmt.Errorf("Could not parse template: %s", err)
			}
			f, err := os.Create(path)
			if err != nil {
				return fmt.Errorf("Could not create file: %s", err)
			}
			defer func() {
				helper.LogE("Could not close file", f.Close())
			}()

			if err := tpl.Execute(f, cfg); err != nil {
				return fmt.Errorf("Could not execute template: %s", err)
			}

			return nil
		},
	)
}
