package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"act-cli/internal/kit"
)

type initOptions struct {
	target             string
	explicitKit        string
	overwrite          bool
	interactive        bool
	assumeYes          bool
	autoApproveDeps    bool
	scope              string
	kitRepo            string
	installSkillsDeps  bool
	installSkillsIsSet bool
}

func RunInit(args []string, w io.Writer) error {
	opts, err := parseInitArgs(args)
	if err != nil {
		return err
	}

	useInteractive := opts.interactive && isTerminalInput()
	if useInteractive {
		var wizardErr error
		opts, wizardErr = runInitWizard(os.Stdin, w, opts)
		if wizardErr != nil {
			return wizardErr
		}
	}

	if opts.scope == kit.ScopeGlobal {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return homeErr
		}
		opts.target = home
	}
	if err := applyInitEnv(opts.target); err != nil {
		return err
	}

	preparedKit, err := kit.PrepareKit(kit.PrepareOptions{
		ExplicitKitPath: opts.explicitKit,
		KitRepo:         opts.kitRepo,
		Out:             w,
	})
	if err != nil {
		return err
	}
	if preparedKit.TempDir != "" {
		defer os.RemoveAll(preparedKit.TempDir)
	}

	abs, err := filepath.Abs(opts.target)
	if err != nil {
		return err
	}
	if _, err := os.Stat(abs); err != nil {
		return err
	}
	if err := kit.InstallRuntime(kit.InstallOptions{
		TargetDir: abs,
		KitPath:   preparedKit.Path,
		Overwrite: opts.overwrite,
		Scope:     opts.scope,
	}); err != nil {
		return err
	}
	if opts.installSkillsDeps {
		fmt.Fprintln(w, "Installing skills dependencies...")
		approveDeps := opts.assumeYes || opts.autoApproveDeps
		if err := kit.InstallSkillsDependencies(abs, approveDeps, w); err != nil {
			return err
		}
	}
	if opts.scope == kit.ScopeGlobal {
		fmt.Fprintf(w, "Initialized act-kit runtime globally in %s\\.claude\n", abs)
	} else {
		fmt.Fprintf(w, "Initialized act-kit runtime in %s\n", abs)
	}
	return nil
}

func parseInitArgs(args []string) (initOptions, error) {
	opts := initOptions{
		target:      ".",
		interactive: true,
		scope:       kit.ScopeProject,
		kitRepo:     kit.DefaultGitHubKitRepo,
	}

	seenTarget := false
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--kit":
			if i+1 >= len(args) {
				return initOptions{}, fmt.Errorf("--kit requires a value")
			}
			opts.explicitKit = args[i+1]
			i++
		case "--force":
			opts.overwrite = true
		case "--yes", "-y":
			opts.assumeYes = true
			opts.interactive = false
		case "--interactive", "-i":
			opts.interactive = true
			opts.assumeYes = false
		case "--non-interactive":
			opts.interactive = false
		case "--global", "-g":
			opts.scope = kit.ScopeGlobal
		case "--project":
			opts.scope = kit.ScopeProject
		case "--kit-repo":
			if i+1 >= len(args) {
				return initOptions{}, fmt.Errorf("--kit-repo requires a value (owner/repo)")
			}
			opts.kitRepo = strings.TrimSpace(args[i+1])
			i++
		case "--install-skills-deps":
			opts.installSkillsDeps = true
			opts.installSkillsIsSet = true
		case "--skip-skills-deps":
			opts.installSkillsDeps = false
			opts.installSkillsIsSet = true
		case "--scope":
			if i+1 >= len(args) {
				return initOptions{}, fmt.Errorf("--scope requires a value (project|global)")
			}
			scope, scopeErr := normalizeScope(args[i+1])
			if scopeErr != nil {
				return initOptions{}, scopeErr
			}
			opts.scope = scope
			i++
		default:
			if strings.HasPrefix(arg, "-") {
				return initOptions{}, fmt.Errorf("unknown option %q", arg)
			}
			if seenTarget {
				return initOptions{}, fmt.Errorf("unexpected argument %q", arg)
			}
			opts.target = arg
			seenTarget = true
		}
	}

	return opts, nil
}

func parseKitAndForce(args []string) (string, bool, error) {
	opts, err := parseInitArgs(args)
	if err != nil {
		return "", false, err
	}
	if opts.target != "." {
		return "", false, fmt.Errorf("unexpected argument %q", opts.target)
	}
	kitPath, err := kit.ResolveKitPath(opts.explicitKit)
	if err != nil {
		return "", false, err
	}
	return kitPath, opts.overwrite, nil
}

func normalizeScope(scope string) (string, error) {
	switch strings.ToLower(strings.TrimSpace(scope)) {
	case "project", "local":
		return kit.ScopeProject, nil
	case "global":
		return kit.ScopeGlobal, nil
	default:
		return "", fmt.Errorf("invalid scope %q; expected project or global", scope)
	}
}

func isTerminalInput() bool {
	info, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (info.Mode() & os.ModeCharDevice) != 0
}
