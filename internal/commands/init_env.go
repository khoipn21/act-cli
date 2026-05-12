package commands

import (
	"bufio"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func applyInitEnv(targetDir string) error {
	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return err
	}
	localEnv := filepath.Join(absTarget, ".env")
	if fileExists(localEnv) {
		if _, err := loadEnvFromFile(localEnv, true); err != nil {
			return err
		}
		return nil
	}

	for _, f := range globalEnvCandidates() {
		if !fileExists(f) {
			continue
		}
		if _, err := loadEnvFromFile(f, false); err != nil {
			return err
		}
	}
	return nil
}

func globalEnvCandidates() []string {
	home, err := os.UserHomeDir()
	if err != nil || strings.TrimSpace(home) == "" {
		return nil
	}
	candidates := []string{
		filepath.Join(home, ".env"),
		filepath.Join(home, ".bashrc"),
		filepath.Join(home, ".bash_profile"),
		filepath.Join(home, ".zshrc"),
		filepath.Join(home, ".profile"),
	}
	if runtime.GOOS == "windows" {
		candidates = append(candidates,
			filepath.Join(home, "Documents", "PowerShell", "Microsoft.PowerShell_profile.ps1"),
			filepath.Join(home, ".config", "powershell", "Microsoft.PowerShell_profile.ps1"),
		)
	}
	return candidates
}

func loadEnvFromFile(path string, overwrite bool) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	loaded := 0
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		k, v, ok := parseEnvAssignment(sc.Text())
		if !ok {
			continue
		}
		if !overwrite && strings.TrimSpace(os.Getenv(k)) != "" {
			continue
		}
		if err := os.Setenv(k, v); err != nil {
			return loaded, err
		}
		loaded++
	}
	if err := sc.Err(); err != nil {
		return loaded, err
	}
	return loaded, nil
}

func parseEnvAssignment(line string) (string, string, bool) {
	s := strings.TrimSpace(line)
	if s == "" || strings.HasPrefix(s, "#") {
		return "", "", false
	}
	s = strings.TrimPrefix(s, "export ")
	if strings.HasPrefix(strings.ToLower(s), "$env:") {
		s = s[len("$env:"):]
	}

	sep := "="
	if strings.Contains(s, "=") {
		sep = "="
	} else {
		return "", "", false
	}
	idx := strings.Index(s, sep)
	if idx <= 0 {
		return "", "", false
	}
	key := strings.TrimSpace(s[:idx])
	val := strings.TrimSpace(s[idx+1:])
	if key == "" {
		return "", "", false
	}
	val = strings.TrimPrefix(val, " ")
	val = strings.Trim(val, `"'`)
	return key, val, true
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
