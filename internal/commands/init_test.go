package commands

import "testing"

func TestParseInitArgsDefaults(t *testing.T) {
	opts, err := parseInitArgs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.target != "." {
		t.Fatalf("expected default target '.', got %q", opts.target)
	}
	if opts.explicitKit != "" {
		t.Fatalf("expected empty kit path, got %q", opts.explicitKit)
	}
	if opts.overwrite {
		t.Fatal("expected overwrite=false by default")
	}
	if !opts.interactive {
		t.Fatal("expected interactive=true by default")
	}
	if opts.assumeYes {
		t.Fatal("expected assumeYes=false by default")
	}
	if opts.scope != "project" {
		t.Fatalf("expected default scope=project, got %q", opts.scope)
	}
	if opts.kitRepo == "" {
		t.Fatal("expected default kit repo to be set")
	}
	if opts.installSkillsDeps {
		t.Fatal("expected installSkillsDeps=false by default")
	}
}

func TestParseInitArgsWithFlagsAndTarget(t *testing.T) {
	opts, err := parseInitArgs([]string{"./proj", "--kit", "C:\\kit", "--force", "--interactive"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.target != "./proj" {
		t.Fatalf("expected target ./proj, got %q", opts.target)
	}
	if opts.explicitKit != "C:\\kit" {
		t.Fatalf("expected kit path C:\\kit, got %q", opts.explicitKit)
	}
	if !opts.overwrite {
		t.Fatal("expected overwrite=true")
	}
	if !opts.interactive {
		t.Fatal("expected interactive=true")
	}
	if opts.scope != "project" {
		t.Fatalf("expected scope=project, got %q", opts.scope)
	}
}

func TestParseInitArgsScopeAndNonInteractive(t *testing.T) {
	opts, err := parseInitArgs([]string{"--scope", "global", "--non-interactive", "--yes"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.scope != "global" {
		t.Fatalf("expected scope=global, got %q", opts.scope)
	}
	if opts.interactive {
		t.Fatal("expected interactive=false with --non-interactive")
	}
	if !opts.assumeYes {
		t.Fatal("expected assumeYes=true with --yes")
	}
}

func TestParseInitArgsGlobalAliasAndKitRepo(t *testing.T) {
	opts, err := parseInitArgs([]string{"-g", "--kit-repo", "owner/repo"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.scope != "global" {
		t.Fatalf("expected scope=global, got %q", opts.scope)
	}
	if opts.kitRepo != "owner/repo" {
		t.Fatalf("expected kitRepo owner/repo, got %q", opts.kitRepo)
	}
}

func TestParseInitArgsInstallSkillsFlags(t *testing.T) {
	opts, err := parseInitArgs([]string{"--install-skills-deps"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !opts.installSkillsDeps || !opts.installSkillsIsSet {
		t.Fatal("expected install-skills-deps to set installSkillsDeps")
	}

	opts, err = parseInitArgs([]string{"--skip-skills-deps"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.installSkillsDeps || !opts.installSkillsIsSet {
		t.Fatal("expected skip-skills-deps to clear installSkillsDeps")
	}
}

func TestParseInitArgsInvalidScope(t *testing.T) {
	if _, err := parseInitArgs([]string{"--scope", "team"}); err == nil {
		t.Fatal("expected error for invalid scope")
	}
}

func TestParseInitArgsUnknownFlag(t *testing.T) {
	if _, err := parseInitArgs([]string{"--bogus"}); err == nil {
		t.Fatal("expected error for unknown flag")
	}
}

func TestParseInitArgsMultipleTargets(t *testing.T) {
	if _, err := parseInitArgs([]string{"a", "b"}); err == nil {
		t.Fatal("expected error for multiple target args")
	}
}
