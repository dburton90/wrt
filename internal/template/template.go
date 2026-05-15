// Package template provides file-copy-with-substitution for task scaffolding.
//
// Apply walks a source directory, mirrors its structure into a destination
// directory, and copies each file with simple {key} → value substitution
// applied to file contents. Filenames and directory names are not substituted.
// Destination files that already exist are skipped, never overwritten.
package template

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Apply walks srcDir and copies every file into dstDir, performing variable
// substitution on file contents. Variables in `vars` are referenced as
// `{key}` in template files. Destination files that already exist are left
// untouched. Empty directories in srcDir are mirrored.
func Apply(srcDir, dstDir string, vars map[string]string) error {
	pairs := make([]string, 0, len(vars)*2)
	for k, v := range vars {
		pairs = append(pairs, "{"+k+"}", v)
	}
	replacer := strings.NewReplacer(pairs...)

	return filepath.WalkDir(srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}
		dst := filepath.Join(dstDir, rel)

		if d.IsDir() {
			return os.MkdirAll(dst, 0o755)
		}

		if _, err := os.Stat(dst); err == nil {
			return nil
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading template %s: %w", path, err)
		}
		if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
			return fmt.Errorf("creating dir for %s: %w", dst, err)
		}
		return os.WriteFile(dst, []byte(replacer.Replace(string(data))), 0o644)
	})
}
