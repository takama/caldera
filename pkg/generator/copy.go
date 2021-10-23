package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/takama/caldera/pkg/helper"
)

// copyTemplates copies templates files from source to destination directory.
func copyTemplates(src, dst string) error {
	if err := filepath.Walk(
		src,
		func(srcPath string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("failed to scan template directory: %w", err)
			}
			dstPath := strings.Replace(srcPath, src, dst, 1)
			if info.IsDir() {
				fi, err := os.Stat(srcPath)
				if err != nil {
					return fmt.Errorf("could not get source directory info: %w", err)
				}
				if err := os.MkdirAll(dstPath, fi.Mode()); err != nil {
					return fmt.Errorf("could not create destination directory: %w", err)
				}
			} else if err := copyFile(srcPath, dstPath, info); err != nil {
				return fmt.Errorf("could not copy file: %w", err)
			}

			return nil
		},
	); err != nil {
		return fmt.Errorf("failed to copy template files %w", err)
	}

	return nil
}

func copyFile(src, dst string, info os.FileInfo) error {
	srcF, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file: %w", err)
	}

	defer func() {
		helper.LogE("could not close file", srcF.Close())
	}()

	dstF, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("could not create destination file: %w", err)
	}

	defer func() {
		helper.LogE("could not close file", dstF.Close())
	}()

	if _, err = io.Copy(dstF, srcF); err != nil {
		return fmt.Errorf("could not copy file: %w", err)
	}

	if err := os.Chmod(dst, info.Mode()); err != nil {
		return fmt.Errorf("failed to assing file mode %w", err)
	}

	return nil
}
