package generator

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/takama/caldera/pkg/config"
	"github.com/takama/caldera/pkg/helper"
)

const randomStringLength = 32
const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// render executes templates in service directory with configured data.
func render(cfg *config.Config) error {
	if err := filepath.Walk(
		cfg.Directories.Service,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("failed to scan service directory: %w", err)
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
							return strings.ToUpper(strings.ReplaceAll(str, "-", "_"))
						},
						"currentYear": func() int {
							return time.Now().Year()
						},
						"randStr": func() string {
							str, err := randomString(randomStringLength, chars)
							if err != nil {
								panic(err)
							}

							return str
						},
					},
				).ParseFiles(path)
			if err != nil {
				return fmt.Errorf("could not parse template: %w", err)
			}
			f, err := os.Create(path)
			if err != nil {
				return fmt.Errorf("could not create file: %w", err)
			}
			defer func() {
				helper.LogE("could not close file", f.Close())
			}()

			if err := tpl.Execute(f, cfg); err != nil {
				return fmt.Errorf("could not execute template: %w", err)
			}

			return nil
		},
	); err != nil {
		return fmt.Errorf("failed to render template files %w", err)
	}

	return nil
}

// randomString generates a random string with the requested length from the source of chars.
func randomString(length int, source string) (string, error) {
	bytes := make([]byte, length)

	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("could not generate random string: %w", err)
	}

	for i, b := range bytes {
		bytes[i] = source[b%byte(len(source))]
	}

	return string(bytes), nil
}
