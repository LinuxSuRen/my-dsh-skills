package skills

import (
	"bytes"
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"
)

// Meta is the frontmatter of a SKILL.md file.
type Meta struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	WhenToUse   string `yaml:"whenToUse"`
}

// ParseFrontmatter extracts the YAML frontmatter delimited by leading
// `---` markers from a SKILL.md document.
func ParseFrontmatter(data []byte) (*Meta, error) {
	const delim = "---\n"
	if !bytes.HasPrefix(data, []byte(delim)) {
		return nil, errors.New("missing leading frontmatter delimiter")
	}
	end := bytes.Index(data[len(delim):], []byte(delim))
	if end < 0 {
		return nil, errors.New("missing closing frontmatter delimiter")
	}
	var meta Meta
	if err := yaml.Unmarshal(data[len(delim):len(delim)+end], &meta); err != nil {
		return nil, fmt.Errorf("parse frontmatter: %w", err)
	}
	if meta.Name == "" || meta.Description == "" {
		return nil, errors.New("frontmatter requires name and description")
	}
	return &meta, nil
}
