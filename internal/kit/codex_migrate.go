package kit

import (
	"fmt"
	"os"
	"path/filepath"

	"act-cli/internal/fsutil"
)

type CodexMigrateOptions struct {
	TargetDir string
	Scope     string
	Overwrite bool
}

type CodexMigrateResult struct {
	SourceClaudeDir string
	AgentsFile      string
	SkillsDir       string
	SubagentsDir    string
	RulesDir        string
	HooksDir        string
	HooksJSON       string
	ConfigFile      string
	AgentsCount     int
	HooksCount      int
	WrappersCount   int
}

func MigrateToCodex(opt CodexMigrateOptions) (CodexMigrateResult, error) {
	scope := opt.Scope
	if scope == "" {
		scope = ScopeProject
	}
	if scope != ScopeProject && scope != ScopeGlobal {
		return CodexMigrateResult{}, fmt.Errorf("invalid scope %q; expected project or global", scope)
	}

	sourceClaudeDir := filepath.Join(opt.TargetDir, ".claude")
	if _, err := os.Stat(sourceClaudeDir); err != nil {
		return CodexMigrateResult{}, fmt.Errorf("source runtime not found at %s", sourceClaudeDir)
	}

	result := CodexMigrateResult{
		SourceClaudeDir: sourceClaudeDir,
		AgentsFile:      filepath.Join(opt.TargetDir, "AGENTS.md"),
		SkillsDir:       filepath.Join(opt.TargetDir, ".agents", "skills"),
		SubagentsDir:    filepath.Join(opt.TargetDir, ".codex", "agents"),
		RulesDir:        filepath.Join(opt.TargetDir, ".codex", "rules"),
		HooksDir:        filepath.Join(opt.TargetDir, ".codex", "hooks"),
		HooksJSON:       filepath.Join(opt.TargetDir, ".codex", "hooks.json"),
		ConfigFile:      filepath.Join(opt.TargetDir, ".codex", "config.toml"),
	}
	for _, dir := range []string{result.SkillsDir, result.SubagentsDir, result.RulesDir, result.HooksDir} {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return CodexMigrateResult{}, err
		}
	}

	agentsSrc := filepath.Join(sourceClaudeDir, "AGENTS.md")
	if _, err := os.Stat(agentsSrc); err != nil {
		agentsSrc = filepath.Join(sourceClaudeDir, "CLAUDE.md")
	}
	if err := copyFileIfExists(agentsSrc, result.AgentsFile, opt.Overwrite); err != nil {
		return CodexMigrateResult{}, err
	}

	if err := fsutil.CopyTree(filepath.Join(sourceClaudeDir, "skills"), result.SkillsDir, fsutil.CopyOptions{Overwrite: opt.Overwrite}); err != nil {
		return CodexMigrateResult{}, fmt.Errorf("copy skills: %w", err)
	}
	if err := fsutil.CopyTree(filepath.Join(sourceClaudeDir, "rules"), result.RulesDir, fsutil.CopyOptions{Overwrite: opt.Overwrite}); err != nil {
		return CodexMigrateResult{}, fmt.Errorf("copy rules: %w", err)
	}
	if err := fsutil.CopyTree(filepath.Join(sourceClaudeDir, "hooks"), result.HooksDir, fsutil.CopyOptions{Overwrite: opt.Overwrite}); err != nil {
		return CodexMigrateResult{}, fmt.Errorf("copy hooks: %w", err)
	}

	agentEntries, agentCount, err := writeCodexAgentTomlFiles(filepath.Join(sourceClaudeDir, "agents"), result.SubagentsDir, opt.Overwrite)
	if err != nil {
		return CodexMigrateResult{}, err
	}
	result.AgentsCount = agentCount

	hooksCount, wrappersCount, err := migrateClaudeHooksToCodex(
		filepath.Join(sourceClaudeDir, "settings.json"),
		result.HooksDir,
		result.HooksJSON,
		opt.Overwrite,
	)
	if err != nil {
		return CodexMigrateResult{}, err
	}
	result.HooksCount = hooksCount
	result.WrappersCount = wrappersCount

	if err := mergeCodexConfigToml(result.ConfigFile, agentEntries, opt.Overwrite); err != nil {
		return CodexMigrateResult{}, err
	}
	return result, nil
}
