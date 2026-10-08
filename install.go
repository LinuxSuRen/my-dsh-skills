package skills

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// Install copies the given skills (all when empty) from the embedded
// filesystem into target, overwriting previous versions. It returns the
// names of the installed skills.
func Install(target string, names []string) ([]string, error) {
	if len(names) == 0 {
		var err error
		names, err = Names()
		if err != nil {
			return nil, err
		}
	}
	if err := os.MkdirAll(target, 0o755); err != nil {
		return nil, fmt.Errorf("create target directory: %w", err)
	}
	for _, name := range names {
		if err := copySkill(fsys, name, filepath.Join(target, name)); err != nil {
			return nil, fmt.Errorf("install %s: %w", name, err)
		}
	}
	return names, nil
}

// Uninstall removes the given skills (all when empty) from target.
func Uninstall(target string, names []string) ([]string, error) {
	if len(names) == 0 {
		var err error
		names, err = Names()
		if err != nil {
			return nil, err
		}
	}
	for _, name := range names {
		dir := filepath.Join(target, name)
		if err := os.RemoveAll(dir); err != nil {
			return nil, fmt.Errorf("uninstall %s: %w", name, err)
		}
	}
	return names, nil
}

func copySkill(fsys fs.FS, name, dst string) error {
	if err := os.RemoveAll(dst); err != nil {
		return fmt.Errorf("clean previous install: %w", err)
	}
	return fs.WalkDir(fsys, name, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(name, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := fsys.Open(path)
		if err != nil {
			return err
		}
		defer func() { _ = data.Close() }()
		out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
		if err != nil {
			return err
		}
		_, copyErr := io.Copy(out, data)
		closeErr := out.Close()
		if copyErr != nil {
			return copyErr
		}
		return closeErr
	})
}
