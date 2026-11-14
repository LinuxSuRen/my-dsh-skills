// Package skills embeds the skill directories of this repository and
// provides APIs to enumerate and install them for AI coding tools.
package skills

import (
	"embed"
	"io/fs"
)

// Skill directories live under the skills/ subtree and are embedded via
// the directory pattern below: adding a new skill requires no
// registration in this file, only a directory with a valid SKILL.md
// under skills/.
//
//go:embed skills
var fsys embed.FS

// root is the embedded skills/ subtree. fs.Sub can only fail when the
// directory is missing, which the embed directive above precludes.
var root fs.FS = mustSub(fsys, "skills")

func mustSub(f fs.FS, dir string) fs.FS {
	sub, err := fs.Sub(f, dir)
	if err != nil {
		panic(err)
	}
	return sub
}

// FS returns the read-only filesystem with all embedded skills, rooted
// at skills/ so that each child directory is one skill.
func FS() fs.FS { return root }
