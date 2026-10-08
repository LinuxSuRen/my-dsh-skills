package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseFrontmatter(t *testing.T) {
	meta, err := ParseFrontmatter([]byte("---\nname: demo\ndescription: A demo skill.\nwhenToUse: Always.\n---\nbody"))
	if err != nil {
		t.Fatalf("parse: %v", err)
	}
	if meta.Name != "demo" || meta.Description != "A demo skill." || meta.WhenToUse != "Always." {
		t.Fatalf("unexpected meta: %+v", meta)
	}
}

func TestParseFrontmatterErrors(t *testing.T) {
	cases := [][]byte{
		[]byte("no frontmatter"),
		[]byte("---\nname: demo\n"),
		[]byte("---\ndescription: missing name\n---\n"),
	}
	for i, c := range cases {
		if _, err := ParseFrontmatter(c); err == nil {
			t.Fatalf("case %d: expected error", i)
		}
	}
}

func TestAllMatchesEmbeddedDirs(t *testing.T) {
	metas, err := All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(metas) == 0 {
		t.Fatal("no embedded skills found")
	}
	for _, m := range metas {
		if _, err := FS().Open(m.Name + "/SKILL.md"); err != nil {
			t.Errorf("skill %s: %v", m.Name, err)
		}
	}
}

func TestTargetDir(t *testing.T) {
	for _, tool := range []string{ToolAgents, ToolDSH, ToolOpenCode, ToolCodex} {
		dir, err := TargetDir(tool)
		if err != nil {
			t.Fatalf("tool %s: %v", tool, err)
		}
		if !filepath.IsAbs(dir) {
			t.Errorf("tool %s: not absolute: %s", tool, dir)
		}
	}
	if _, err := TargetDir("nope"); err == nil {
		t.Fatal("expected error for unknown tool")
	}
}

func TestInstallAndUninstall(t *testing.T) {
	target := t.TempDir()

	installed, err := Install(target, nil)
	if err != nil {
		t.Fatalf("install: %v", err)
	}
	if len(installed) == 0 {
		t.Fatal("installed nothing")
	}
	for _, name := range installed {
		if _, err := os.Stat(filepath.Join(target, name, "SKILL.md")); err != nil {
			t.Errorf("skill %s not installed: %v", name, err)
		}
	}

	// Install again to a subdir skill with extra files must be idempotent.
	if _, err := Install(target, []string{"skill-refresh"}); err != nil {
		t.Fatalf("reinstall: %v", err)
	}
	if _, err := os.Stat(filepath.Join(target, "skill-refresh", "state.md")); err != nil {
		t.Errorf("extra skill file missing: %v", err)
	}

	removed, err := Uninstall(target, []string{"skill-refresh"})
	if err != nil {
		t.Fatalf("uninstall: %v", err)
	}
	if len(removed) != 1 {
		t.Fatalf("unexpected removed: %v", removed)
	}
	if _, err := os.Stat(filepath.Join(target, "skill-refresh")); !os.IsNotExist(err) {
		t.Errorf("skill-refresh still present: %v", err)
	}
}
