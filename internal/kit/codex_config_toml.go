package kit

import (
	"os"
	"path/filepath"
	"strings"
)

const (
	codexAgentsBlockStart = "# BEGIN ACT-CODEX-AGENTS"
	codexAgentsBlockEnd   = "# END ACT-CODEX-AGENTS"
)

func mergeCodexConfigToml(path string, agentEntries []string, overwrite bool) error {
	existing := ""
	if b, err := os.ReadFile(path); err == nil {
		existing = string(b)
	} else if !os.IsNotExist(err) {
		return err
	}
	if existing == "" {
		existing = "model = \"gpt-5\"\nmodel_reasoning_effort = \"high\"\n"
	}
	if !overwrite {
		// In non-overwrite mode we still keep config consistent by inserting
		// managed blocks only when missing.
		existing = ensureHooksFeatureFlag(existing)
		if strings.Contains(existing, codexAgentsBlockStart) {
			return writeText(path, existing)
		}
	}

	managedBody := codexAgentsBlockStart + "\n"
	if len(agentEntries) > 0 {
		managedBody += strings.Join(agentEntries, "\n\n") + "\n"
	}
	managedBody += codexAgentsBlockEnd

	updated := replaceOrAppendManagedBlock(existing, managedBody)
	updated = ensureHooksFeatureFlag(updated)
	return writeText(path, updated)
}

func replaceOrAppendManagedBlock(content, managedBlock string) string {
	start := strings.Index(content, codexAgentsBlockStart)
	end := strings.Index(content, codexAgentsBlockEnd)
	if start >= 0 && end > start {
		end += len(codexAgentsBlockEnd)
		return strings.TrimSpace(content[:start]) + "\n\n" + managedBlock + "\n"
	}
	content = strings.TrimSpace(content)
	if content == "" {
		return managedBlock + "\n"
	}
	return content + "\n\n" + managedBlock + "\n"
}

func ensureHooksFeatureFlag(content string) string {
	lines := strings.Split(strings.ReplaceAll(content, "\r\n", "\n"), "\n")
	filtered := make([]string, 0, len(lines))
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "codex_hooks") {
			continue
		}
		filtered = append(filtered, line)
	}
	lines = filtered

	featuresStart, featuresEnd := -1, len(lines)
	for i, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "[features]" {
			featuresStart = i
			continue
		}
		if featuresStart >= 0 && strings.HasPrefix(trim, "[") && strings.HasSuffix(trim, "]") {
			featuresEnd = i
			break
		}
	}

	if featuresStart == -1 {
		content = strings.TrimSpace(strings.Join(lines, "\n"))
		if content == "" {
			return "[features]\nhooks = true\n"
		}
		return content + "\n\n[features]\nhooks = true\n"
	}

	foundHooks := false
	for i := featuresStart + 1; i < featuresEnd; i++ {
		trim := strings.TrimSpace(lines[i])
		if strings.HasPrefix(trim, "hooks") && strings.Contains(trim, "=") {
			lines[i] = "hooks = true"
			foundHooks = true
		}
	}
	if !foundHooks {
		before := append([]string{}, lines[:featuresEnd]...)
		after := append([]string{}, lines[featuresEnd:]...)
		before = append(before, "hooks = true")
		lines = append(before, after...)
	}

	return strings.TrimSpace(strings.Join(lines, "\n")) + "\n"
}

func writeText(path, content string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0o644)
}
