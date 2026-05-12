package commands

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseEnvAssignment(t *testing.T) {
	tests := []struct {
		in   string
		key  string
		val  string
		want bool
	}{
		{"FOO=bar", "FOO", "bar", true},
		{"export HELLO=world", "HELLO", "world", true},
		{"$env:TOKEN = \"abc\"", "TOKEN", "abc", true},
		{"# comment", "", "", false},
		{"", "", "", false},
	}
	for _, tc := range tests {
		k, v, ok := parseEnvAssignment(tc.in)
		if ok != tc.want || k != tc.key || v != tc.val {
			t.Fatalf("parseEnvAssignment(%q) = (%q,%q,%v), want (%q,%q,%v)", tc.in, k, v, ok, tc.key, tc.val, tc.want)
		}
	}
}

func TestApplyInitEnvPrefersLocalEnv(t *testing.T) {
	dir := t.TempDir()
	localEnv := filepath.Join(dir, ".env")
	if err := os.WriteFile(localEnv, []byte("ACT_GITHUB_TOKEN=local-token\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("ACT_GITHUB_TOKEN", "")
	if err := applyInitEnv(dir); err != nil {
		t.Fatalf("applyInitEnv error: %v", err)
	}
	if got := os.Getenv("ACT_GITHUB_TOKEN"); got != "local-token" {
		t.Fatalf("expected local token, got %q", got)
	}
}
