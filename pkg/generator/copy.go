package generator

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/takama/caldera/pkg/helper"
)

// copyTemplates copies templates files from source to destination directory
func copyTemplates(src, dst string) error {
	return filepath.Walk(
		src,
		func(srcPath string, info os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("failed to scan template directory: %s", err)
			}
			dstPath := strings.Replace(srcPath, src, dst, 1)
			if info.IsDir() {
				fi, err := os.Stat(srcPath)
				if err != nil {
					return fmt.Errorf("could not get source directory info: %s", err)
				}
				if err := os.MkdirAll(dstPath, fi.Mode()); err != nil {
					return fmt.Errorf("could not create destination directory: %s", err)
				}
			} else if err := copyFile(srcPath, dstPath, info); err != nil {
				return fmt.Errorf("could not copy file: %s", err)
			}
			return nil
		},
	)
}

func copyFile(src, dst string, info os.FileInfo) error {
	srcF, err := os.Open(src) // nolint: gosec
	if err != nil {
		return fmt.Errorf("could not open source file: %s", err)
	}
	defer func() {
		helper.LogE("could not close file", srcF.Close())
	}()

	dstF, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("could not create destination file: %s", err)
	}
	defer func() {
		helper.LogE("could not close file", dstF.Close())
	}()

	if _, err = io.Copy(dstF, srcF); err != nil {
		return fmt.Errorf("could not copy file: %s", err)
	}
	return os.Chmod(dst, info.Mode())
}
