package kit

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveKitPathWithExplicitValidPath(t *testing.T) {
	root := createTestKitRoot(t)

	got, err := ResolveKitPath(root)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != root {
		t.Fatalf("expected %q, got %q", root, got)
	}
}

func TestResolveKitPathWithExplicitInvalidPath(t *testing.T) {
	root := t.TempDir()

	_, err := ResolveKitPath(root)
	if err == nil {
		t.Fatal("expected error for invalid explicit kit path")
	}
	if !strings.Contains(err.Error(), "missing claude/metadata.json") {
		t.Fatalf("expected missing metadata error, got %v", err)
	}
}

func TestResolveKitPathAutodetectSkipsInvalidCandidate(t *testing.T) {
	parent := t.TempDir()
	cwd := filepath.Join(parent, "workspace")
	if err := os.MkdirAll(cwd, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(parent, "act-kit"), 0o755); err != nil {
		t.Fatal(err)
	}
	valid := filepath.Join(cwd, "act-kit")
	createTestKitRootAt(t, valid)

	previousWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(previousWD)
	})
	if err := os.Chdir(cwd); err != nil {
		t.Fatal(err)
	}

	got, err := ResolveKitPath("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != valid {
		t.Fatalf("expected autodetected kit %q, got %q", valid, got)
	}
}

func TestResolveKitPathAutodetectMissing(t *testing.T) {
	cwd := t.TempDir()
	previousWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(previousWD)
	})
	if err := os.Chdir(cwd); err != nil {
		t.Fatal(err)
	}

	_, err = ResolveKitPath("")
	if err == nil {
		t.Fatal("expected missing local kit error")
	}
	if !strings.Contains(err.Error(), "local kit source not found") {
		t.Fatalf("expected missing local kit error, got %v", err)
	}
}

func createTestKitRoot(t *testing.T) string {
	t.Helper()
	root := filepath.Join(t.TempDir(), "act-kit")
	createTestKitRootAt(t, root)
	return root
}

func createTestKitRootAt(t *testing.T, root string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Join(root, "claude"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "claude", "metadata.json"), []byte(`{"name":"act-kit","version":"test"}`), 0o644); err != nil {
		t.Fatal(err)
	}
}
