package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunInitWizardHidesDefaultsAndAutoApprovesDeps(t *testing.T) {
	t.Setenv("NO_COLOR", "1")

	targetDir := filepath.Join(t.TempDir(), "project")
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		t.Fatalf("mkdir target: %v", err)
	}

	in := strings.NewReader("\n\ny\ny\n")
	var out bytes.Buffer
	opts := initOptions{
		target:      targetDir,
		scope:       "project",
		kitRepo:     "khoipn21/act-kit",
		interactive: true,
	}

	got, err := runInitWizard(in, &out, opts)
	if err != nil {
		t.Fatalf("runInitWizard error: %v", err)
	}
	if !got.installSkillsDeps {
		t.Fatal("expected installSkillsDeps=true")
	}
	if !got.autoApproveDeps {
		t.Fatal("expected autoApproveDeps=true after wizard confirmation")
	}

	text := out.String()
	for _, forbidden := range []string{
		"[Defaults]",
		"Target directory:",
		"Kit source:",
		"?? Overwrite existing files?",
		"?? Install skills dependencies now?",
	} {
		if strings.Contains(text, forbidden) {
			t.Fatalf("wizard output should not contain %q", forbidden)
		}
	}
}

