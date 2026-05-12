package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"act-cli/internal/kit"
)

type migrateOptions struct {
	target      string
	explicitKit string
	overwrite   bool
	to          string
	scope       string
	scopeIsSet  bool
	interactive bool
	assumeYes   bool
}

func RunMigrate(args []string, w io.Writer) error {
	opts, err := parseMigrateArgs(args)
	if err != nil {
		return err
	}

	useInteractive := opts.interactive
	if opts.to == "" {
		if useInteractive {
			opts, err = runMigrateWizardFull(os.Stdin, w, opts)
			if err != nil {
				return err
			}
		} else {
			opts.to = "codex"
		}
	} else if useInteractive && !opts.assumeYes && strings.EqualFold(opts.to, "codex") && !opts.scopeIsSet {
		opts, err = runMigrateWizardScopeOnly(os.Stdin, w, opts)
		if err != nil {
			return err
		}
	}

	if !strings.EqualFold(strings.TrimSpace(opts.to), "codex") {
		return fmt.Errorf("unsupported migration target %q; currently only `codex` is supported", opts.to)
	}

	baseTarget := opts.target
	if opts.scope == kit.ScopeGlobal {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return homeErr
		}
		baseTarget = home
	}

	abs, err := filepath.Abs(baseTarget)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(abs, 0o755); err != nil {
		return err
	}

	if !useInteractive || opts.assumeYes {
		renderMigrateNoPromptSummary(w, opts, abs)
	}

	srcClaude := filepath.Join(abs, ".claude")
	if _, err := os.Stat(srcClaude); os.IsNotExist(err) {
		kitPath, err := kit.ResolveKitPath(opts.explicitKit)
		if err != nil {
			return err
		}
		if err := kit.InstallRuntime(kit.InstallOptions{
			TargetDir: abs,
			KitPath:   kitPath,
			Overwrite: true,
			Scope:     opts.scope,
		}); err != nil {
			return fmt.Errorf("prepare source runtime at %s: %w", srcClaude, err)
		}
	}

	result, err := kit.MigrateToCodex(kit.CodexMigrateOptions{
		TargetDir: abs,
		Scope:     opts.scope,
		Overwrite: opts.overwrite || opts.assumeYes,
	})
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Migrated ACT runtime to Codex (%s scope)\n", opts.scope)
	fmt.Fprintf(w, "  Source: %s\n", result.SourceClaudeDir)
	fmt.Fprintf(w, "  AGENTS: %s\n", result.AgentsFile)
	fmt.Fprintf(w, "  Skills: %s\n", result.SkillsDir)
	fmt.Fprintf(w, "  Agents: %s\n", result.SubagentsDir)
	fmt.Fprintf(w, "  Rules: %s\n", result.RulesDir)
	fmt.Fprintf(w, "  Hooks: %s\n", result.HooksDir)
	fmt.Fprintf(w, "  Hooks JSON: %s\n", result.HooksJSON)
	fmt.Fprintf(w, "  Config: %s\n", result.ConfigFile)
	fmt.Fprintf(w, "  Migrated agents: %d\n", result.AgentsCount)
	fmt.Fprintf(w, "  Migrated hooks: %d (wrappers: %d)\n", result.HooksCount, result.WrappersCount)
	return nil
}

func parseMigrateArgs(args []string) (migrateOptions, error) {
	opts := migrateOptions{
		target:      ".",
		scope:       kit.ScopeProject,
		interactive: true,
	}

	seenTarget := false
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--kit":
			if i+1 >= len(args) {
				return migrateOptions{}, fmt.Errorf("--kit requires a value")
			}
			opts.explicitKit = strings.TrimSpace(args[i+1])
			i++
		case "--force":
			opts.overwrite = true
		case "--to":
			if i+1 >= len(args) {
				return migrateOptions{}, fmt.Errorf("--to requires a value (codex)")
			}
			opts.to = strings.ToLower(strings.TrimSpace(args[i+1]))
			i++
		case "--scope":
			if i+1 >= len(args) {
				return migrateOptions{}, fmt.Errorf("--scope requires a value (project|global)")
			}
			scope, scopeErr := normalizeScope(args[i+1])
			if scopeErr != nil {
				return migrateOptions{}, scopeErr
			}
			opts.scope = scope
			opts.scopeIsSet = true
			i++
		case "--global", "-g":
			opts.scope = kit.ScopeGlobal
			opts.scopeIsSet = true
		case "--project":
			opts.scope = kit.ScopeProject
			opts.scopeIsSet = true
		case "--yes", "-y":
			opts.assumeYes = true
			opts.interactive = false
		case "--interactive", "-i":
			opts.interactive = true
			opts.assumeYes = false
		case "--non-interactive":
			opts.interactive = false
		default:
			if strings.HasPrefix(arg, "-") {
				return migrateOptions{}, fmt.Errorf("unknown option %q", arg)
			}
			if seenTarget {
				return migrateOptions{}, fmt.Errorf("unexpected argument %q", arg)
			}
			opts.target = arg
			seenTarget = true
		}
	}

	return opts, nil
}

func renderMigrateNoPromptSummary(w io.Writer, opts migrateOptions, target string) {
	ui := newMigrateWizardUI(w)
	ui.header()
	ui.section("Target")
	ui.info("Migrate target: codex")
	ui.section("Scope")
	ui.info("Installation scope: " + opts.scope)
	ui.info("Resolved base path: " + target)
	ui.section("Mode")
	ui.info("Non-interactive: defaults applied")
	ui.footer("Starting migration...")
}
