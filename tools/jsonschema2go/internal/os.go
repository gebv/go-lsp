package internal

import (
	"io/fs"
	"path/filepath"
)

func ListFiles(workDir, pattern string) ([]string, error) {
	var a []string
	err := filepath.WalkDir(workDir, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if matched, err := filepath.Match(pattern, d.Name()); err != nil {
			return err
		} else if matched {
			a = append(a, s)
		}
		return nil
	})
	return a, err
}
