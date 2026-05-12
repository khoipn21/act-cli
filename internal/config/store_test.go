package config

import (
	"path/filepath"
	"testing"
)

func TestStoreLoadSave(t *testing.T) {
	tmp := t.TempDir()
	store := &Store{Path: filepath.Join(tmp, "config.json")}

	cfg, err := store.Load()
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(cfg) != 0 {
		t.Fatal("expected empty config")
	}

	cfg["foo"] = "bar"
	if err := store.Save(cfg); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	reloaded, err := store.Load()
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	if reloaded["foo"] != "bar" {
		t.Fatalf("expected foo=bar, got %q", reloaded["foo"])
	}
}
