package generator

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/takama/caldera/pkg/config"
	"github.com/takama/caldera/pkg/helper"
)

// render executes templates in service directory with configured data
func render(cfg *config.Config) error {
	return filepath.Walk(
		cfg.Directories.Service,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("failed to scan service directory: %s", err)
			}
			if info.IsDir() {
				return nil
			}
			name := filepath.Base(path)
			tpl, err := template.New(name).
				Delims("{{[", "]}}").
				Funcs(
					template.FuncMap{
						"toENV": func(str string) string {
							return strings.ToUpper(strings.Replace(str, "-", "_", -1))
						},
						"currentYear": func() int {
							return time.Now().Year()
						},
						"randStr": func() string {
							var ch = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
							b := make([]rune, 32)
							for i := range b {
								b[i] = ch[rand.Intn(len(ch))]
							}
							return string(b)
						},
					},
				).ParseFiles(path)
			if err != nil {
				return fmt.Errorf("could not parse template: %s", err)
			}
			f, err := os.Create(path)
			if err != nil {
				return fmt.Errorf("could not create file: %s", err)
			}
			defer func() {
				helper.LogE("could not close file", f.Close())
			}()

			if err := tpl.Execute(f, cfg); err != nil {
				return fmt.Errorf("could not execute template: %s", err)
			}

			return nil
		},
	)
}
