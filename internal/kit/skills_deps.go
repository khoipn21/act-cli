package kit

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func InstallSkillsDependencies(targetDir string, assumeYes bool, out io.Writer) error {
	skillsDir := filepath.Join(targetDir, ".claude", "skills")
	if _, err := os.Stat(skillsDir); err != nil {
		return fmt.Errorf("skills directory not found: %w", err)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		script := filepath.Join(skillsDir, "install.ps1")
		args := []string{"-NoProfile", "-ExecutionPolicy", "Bypass", "-File", script}
		if assumeYes {
			args = append(args, "-Y")
		}
		cmd = exec.Command("powershell", args...)
	} else {
		script := filepath.Join(skillsDir, "install.sh")
		args := []string{script}
		if assumeYes {
			args = append(args, "--yes")
		}
		cmd = exec.Command("bash", args...)
	}

	cmd.Dir = skillsDir
	cmd.Stdout = out
	cmd.Stderr = out
	cmd.Env = append(os.Environ(), "NON_INTERACTIVE=1")
	if !assumeYes {
		// Let skill installer prompt as needed in interactive flow.
		cmd.Env = os.Environ()
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run skills installer: %w", err)
	}
	return nil
}
