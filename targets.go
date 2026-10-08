package skills

import (
	"fmt"
	"os"
	"path/filepath"
)

// Supported install targets for AI coding tools.
const (
	ToolAgents   = "agents"   // ~/.agents/skills — shared by dsh, codex, opencode
	ToolDSH      = "dsh"      // ~/.dsh/skills
	ToolOpenCode = "opencode" // ~/.config/opencode/skills
	ToolCodex    = "codex"    // ~/.codex/skills
	ToolAll      = "all"      // every target above
)

var toolOrder = []string{ToolAgents, ToolDSH, ToolOpenCode, ToolCodex}

// TargetDir maps a tool name to its user-level skill directory.
func TargetDir(tool string) (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	switch tool {
	case ToolAgents:
		return filepath.Join(home, ".agents", "skills"), nil
	case ToolDSH:
		return filepath.Join(home, ".dsh", "skills"), nil
	case ToolOpenCode:
		if cfg, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok && cfg != "" {
			return filepath.Join(cfg, "opencode", "skills"), nil
		}
		return filepath.Join(home, ".config", "opencode", "skills"), nil
	case ToolCodex:
		return filepath.Join(home, ".codex", "skills"), nil
	default:
		return "", fmt.Errorf("unknown tool %q (supported: %s, %s)", tool, ToolAll, joinQuoted(toolOrder))
	}
}

// TargetDirs resolves one or more target directories for a --tool value.
func TargetDirs(tool string) (map[string]string, error) {
	if tool == ToolAll {
		dirs := make(map[string]string, len(toolOrder))
		for _, t := range toolOrder {
			dir, err := TargetDir(t)
			if err != nil {
				return nil, err
			}
			dirs[t] = dir
		}
		return dirs, nil
	}
	dir, err := TargetDir(tool)
	if err != nil {
		return nil, err
	}
	return map[string]string{tool: dir}, nil
}

func joinQuoted(items []string) string {
	out := ""
	for i, s := range items {
		if i > 0 {
			out += ", "
		}
		out += fmt.Sprintf("%q", s)
	}
	return out
}
