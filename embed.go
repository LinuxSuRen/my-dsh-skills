// Package skills embeds the skill directories of this repository and
// provides APIs to enumerate and install them for AI coding tools.
package skills

import (
	"embed"
	"io/fs"
)

// Skill directories at the repository root. A new skill must be added
// to this list to get embedded into the binary.
//
//go:embed arch-dev-workflow branching-workflow frontend-dev-workflow git-identity open-source-contribution skill-distillation skill-refresh
var fsys embed.FS

// FS returns the read-only filesystem with all embedded skills.
func FS() fs.FS { return fsys }
