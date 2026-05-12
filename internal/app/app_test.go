package app

import (
	"bytes"
	"testing"
)

func TestRunHelp(t *testing.T) {
	t.Setenv("ACT_DISABLE_AUTO_UPDATE_CHECK", "1")
	var out bytes.Buffer
	if err := Run(nil, &out, &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Len() == 0 {
		t.Fatal("expected help output")
	}
}

func TestRunUnknownCommand(t *testing.T) {
	t.Setenv("ACT_DISABLE_AUTO_UPDATE_CHECK", "1")
	var out bytes.Buffer
	err := Run([]string{"unknown"}, &out, &out)
	if err == nil {
		t.Fatal("expected error")
	}
}
