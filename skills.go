package skills

import (
	"fmt"
	"io/fs"
	"sort"
)

// All returns the metadata of every embedded skill, sorted by name.
func All() ([]Meta, error) {
	entries, err := fs.ReadDir(fsys, ".")
	if err != nil {
		return nil, fmt.Errorf("read embedded skills: %w", err)
	}
	var metas []Meta
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		data, err := fs.ReadFile(fsys, e.Name()+"/SKILL.md")
		if err != nil {
			return nil, fmt.Errorf("read %s/SKILL.md: %w", e.Name(), err)
		}
		meta, err := ParseFrontmatter(data)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", e.Name(), err)
		}
		if meta.Name != e.Name() {
			return nil, fmt.Errorf("%s: frontmatter name %q must match the directory name", e.Name(), meta.Name)
		}
		metas = append(metas, *meta)
	}
	sort.Slice(metas, func(i, j int) bool { return metas[i].Name < metas[j].Name })
	return metas, nil
}

// Names returns the names of all embedded skills.
func Names() ([]string, error) {
	metas, err := All()
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(metas))
	for _, m := range metas {
		names = append(names, m.Name)
	}
	return names, nil
}
