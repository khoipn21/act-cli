package kit

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type claudeAgentFile struct {
	Name        string
	Description string
	Model       string
	Tools       string
	Body        string
}

func writeCodexAgentTomlFiles(sourceAgentsDir, targetAgentsDir string, overwrite bool) ([]string, int, error) {
	entries, count := []string{}, 0
	items, err := os.ReadDir(sourceAgentsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return entries, 0, nil
		}
		return nil, 0, err
	}
	for _, item := range items {
		if item.IsDir() || !strings.HasSuffix(strings.ToLower(item.Name()), ".md") {
			continue
		}
		agentPath := filepath.Join(sourceAgentsDir, item.Name())
		agent, err := parseClaudeAgentFile(agentPath)
		if err != nil {
			return nil, 0, err
		}
		slug := toCodexSlug(agent.Name)
		tomlPath := filepath.Join(targetAgentsDir, slug+".toml")
		if _, err := os.Stat(tomlPath); err == nil && !overwrite {
			entries = append(entries, buildCodexConfigEntry(slug, chooseDescription(agent.Description, agent.Name)))
			count++
			continue
		}
		content := buildCodexAgentToml(agent)
		if err := os.WriteFile(tomlPath, []byte(content), 0o644); err != nil {
			return nil, 0, err
		}
		entries = append(entries, buildCodexConfigEntry(slug, chooseDescription(agent.Description, agent.Name)))
		count++
	}
	return entries, count, nil
}

func parseClaudeAgentFile(path string) (claudeAgentFile, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return claudeAgentFile{}, err
	}
	content := strings.ReplaceAll(string(b), "\r\n", "\n")
	agent := claudeAgentFile{Name: strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))}
	if !strings.HasPrefix(content, "---\n") {
		agent.Body = strings.TrimSpace(content)
		return agent, nil
	}
	parts := strings.SplitN(content, "\n---\n", 2)
	if len(parts) != 2 {
		agent.Body = strings.TrimSpace(content)
		return agent, nil
	}
	fm := strings.TrimPrefix(parts[0], "---\n")
	for _, line := range strings.Split(fm, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.Index(line, ":")
		if idx <= 0 {
			continue
		}
		key := strings.TrimSpace(strings.ToLower(line[:idx]))
		val := strings.TrimSpace(line[idx+1:])
		val = strings.Trim(val, "\"'")
		switch key {
		case "name":
			if val != "" {
				agent.Name = val
			}
		case "description":
			agent.Description = val
		case "model":
			agent.Model = val
		case "tools":
			agent.Tools = val
		}
	}
	agent.Body = strings.TrimSpace(parts[1])
	return agent, nil
}

func buildCodexAgentToml(agent claudeAgentFile) string {
	lines := []string{}
	if model := mapModelToCodex(agent.Model); model != "" {
		lines = append(lines, fmt.Sprintf("model = %q", model))
		lines = append(lines, "model_reasoning_effort = \"high\"")
	}
	if mode := deriveSandboxMode(agent.Tools); mode != "" {
		lines = append(lines, fmt.Sprintf("sandbox_mode = %q", mode))
	}
	if len(lines) > 0 {
		lines = append(lines, "")
	}
	lines = append(lines, "developer_instructions = \"\"\"")
	lines = append(lines, escapeTomlMultiline(agent.Body))
	lines = append(lines, "\"\"\"")
	return strings.Join(lines, "\n") + "\n"
}

func mapModelToCodex(model string) string {
	m := strings.ToLower(strings.TrimSpace(model))
	if m == "" || m == "inherit" {
		return "gpt-5"
	}
	if strings.Contains(m, "opus") || strings.Contains(m, "sonnet") || strings.Contains(m, "gpt") {
		return "gpt-5"
	}
	return "gpt-5"
}

func deriveSandboxMode(tools string) string {
	l := strings.ToLower(tools)
	if strings.Contains(l, "bash") || strings.Contains(l, "write") || strings.Contains(l, "edit") || strings.Contains(l, "task") {
		return "workspace-write"
	}
	if strings.Contains(l, "read") || strings.Contains(l, "grep") || strings.Contains(l, "glob") {
		return "read-only"
	}
	return "workspace-write"
}

func escapeTomlMultiline(v string) string {
	return strings.ReplaceAll(v, "\"\"\"", "\\\"\\\"\\\"")
}

func toCodexSlug(name string) string {
	normalized := strings.ToLower(strings.TrimSpace(name))
	re := regexp.MustCompile(`[^a-z0-9]+`)
	slug := strings.Trim(re.ReplaceAllString(normalized, "_"), "_")
	if slug != "" {
		return slug
	}
	h := sha1.Sum([]byte(name))
	return "agent_" + hex.EncodeToString(h[:])[:8]
}

func buildCodexConfigEntry(slug, description string) string {
	return fmt.Sprintf("[agents.%s]\ndescription = %q\nconfig_file = %q", slug, description, "agents/"+slug+".toml")
}

func chooseDescription(description, fallback string) string {
	if strings.TrimSpace(description) != "" {
		return strings.TrimSpace(description)
	}
	return strings.TrimSpace(fallback)
}
