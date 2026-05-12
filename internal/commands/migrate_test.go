package commands

import (
	"testing"

	"act-cli/internal/kit"
)

func TestParseMigrateArgsDefaults(t *testing.T) {
	opts, err := parseMigrateArgs(nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.target != "." {
		t.Fatalf("expected default target '.', got %q", opts.target)
	}
	if opts.scope != kit.ScopeProject {
		t.Fatalf("expected default scope project, got %q", opts.scope)
	}
	if !opts.interactive {
		t.Fatal("expected interactive=true by default")
	}
	if opts.assumeYes {
		t.Fatal("expected assumeYes=false by default")
	}
}

func TestParseMigrateArgsToCodexYes(t *testing.T) {
	opts, err := parseMigrateArgs([]string{"--to", "codex", "-y"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.to != "codex" {
		t.Fatalf("expected to=codex, got %q", opts.to)
	}
	if !opts.assumeYes {
		t.Fatal("expected assumeYes=true")
	}
	if opts.interactive {
		t.Fatal("expected interactive=false when -y is passed")
	}
}

func TestParseMigrateArgsScopeAliases(t *testing.T) {
	opts, err := parseMigrateArgs([]string{"--to", "codex", "-g"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if opts.scope != kit.ScopeGlobal {
		t.Fatalf("expected scope=global, got %q", opts.scope)
	}
	if !opts.scopeIsSet {
		t.Fatal("expected scopeIsSet=true")
	}
}

func TestParseMigrateArgsInvalidScope(t *testing.T) {
	if _, err := parseMigrateArgs([]string{"--scope", "team"}); err == nil {
		t.Fatal("expected invalid scope error")
	}
}

