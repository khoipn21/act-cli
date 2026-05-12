package commands

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintHelpRendersCommandCatalog(t *testing.T) {
	t.Setenv("NO_COLOR", "1")
	var out bytes.Buffer
	PrintHelp(&out)
	s := out.String()
	for _, token := range []string{
		"ACT CLI",
		"[Commands]",
		"init",
		"migrate",
		"update",
		"--to codex",
		"--yes, -y",
		"--check",
	} {
		if !strings.Contains(s, token) {
			t.Fatalf("help output missing token %q", token)
		}
	}
}

