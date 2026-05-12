package kit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMigrateToCodexProjectScope(t *testing.T) {
	root := t.TempDir()
	claudeDir := filepath.Join(root, ".claude")
	if err := os.MkdirAll(filepath.Join(claudeDir, "skills", "demo-skill"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(claudeDir, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(claudeDir, "rules"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(claudeDir, "hooks"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(claudeDir, "AGENTS.md"), []byte("# agents"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(claudeDir, "settings.json"), []byte(`{
  "hooks": {
    "SessionStart": [{"matcher":"startup","hooks":[{"type":"command","command":"node .claude/hooks/demo.cjs"}]}],
    "SubagentStart": [{"hooks":[{"type":"command","command":"node .claude/hooks/skip.cjs"}]}]
  }
}`), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(claudeDir, "skills", "demo-skill", "SKILL.md"), []byte("skill"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(claudeDir, "hooks", "demo.cjs"), []byte("console.log('{}')\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(claudeDir, "agents", "tester.md"), []byte("---\nname: tester\ndescription: test\ntools: Read, Bash\n---\nBody"), 0o644); err != nil {
		t.Fatal(err)
	}

	result, err := MigrateToCodex(CodexMigrateOptions{
		TargetDir: root,
		Scope:     ScopeProject,
		Overwrite: true,
	})
	if err != nil {
		t.Fatalf("unexpected migrate error: %v", err)
	}

	paths := []string{
		result.AgentsFile,
		result.SkillsDir,
		result.SubagentsDir,
		result.RulesDir,
		result.HooksDir,
		result.HooksJSON,
		result.ConfigFile,
		filepath.Join(result.SkillsDir, "demo-skill", "SKILL.md"),
	}
	for _, p := range paths {
		if _, err := os.Stat(p); err != nil {
			t.Fatalf("expected path to exist %s: %v", p, err)
		}
	}
	if result.AgentsCount == 0 {
		t.Fatal("expected migrated agents > 0")
	}
	if result.HooksCount == 0 || result.WrappersCount == 0 {
		t.Fatal("expected hooks and wrappers to be generated")
	}
}
