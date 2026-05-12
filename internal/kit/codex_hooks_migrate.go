package kit

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type claudeSettings struct {
	Hooks map[string][]hookGroup `json:"hooks"`
}

type hookGroup struct {
	Matcher string      `json:"matcher,omitempty"`
	Hooks   []hookEntry `json:"hooks"`
}

type hookEntry struct {
	Type               string `json:"type"`
	Command            string `json:"command"`
	Timeout            int    `json:"timeout,omitempty"`
	PermissionDecision string `json:"permissionDecision,omitempty"`
	Decision           string `json:"decision,omitempty"`
}

func migrateClaudeHooksToCodex(sourceSettingsPath, targetHooksDir, targetHooksJSON string, overwrite bool) (int, int, error) {
	b, err := os.ReadFile(sourceSettingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, 0, nil
		}
		return 0, 0, err
	}

	var settings claudeSettings
	if err := json.Unmarshal(b, &settings); err != nil {
		return 0, 0, err
	}
	filtered := convertHooks(settings.Hooks)
	totalHooks := 0
	wrappers := 0
	for event, groups := range filtered {
		for gi := range groups {
			for hi := range groups[gi].Hooks {
				entry := &groups[gi].Hooks[hi]
				totalHooks++
				wrapperPath := filepath.Join(targetHooksDir, buildHookWrapperName(event, gi, hi, entry.Command))
				originalCommand := rewriteHookCommandPathToCodex(entry.Command)
				if err := writeHookWrapper(wrapperPath, originalCommand, event, overwrite); err != nil {
					return 0, 0, err
				}
				entry.Command = fmt.Sprintf("node %q", wrapperPath)
				wrappers++
			}
		}
		filtered[event] = groups
	}

	out, err := json.MarshalIndent(map[string]any{"hooks": filtered}, "", "  ")
	if err != nil {
		return 0, 0, err
	}
	return totalHooks, wrappers, writeText(targetHooksJSON, string(out)+"\n")
}

func convertHooks(source map[string][]hookGroup) map[string][]hookGroup {
	result := map[string][]hookGroup{}
	for event, groups := range source {
		if !isSupportedHookEvent(event) {
			continue
		}
		kept := make([]hookGroup, 0, len(groups))
		for _, group := range groups {
			normalizedMatcher, ok := normalizeMatcher(event, group.Matcher)
			if !ok {
				continue
			}
			hooks := make([]hookEntry, 0, len(group.Hooks))
			for _, h := range group.Hooks {
				if strings.ToLower(strings.TrimSpace(h.Type)) != "command" || strings.TrimSpace(h.Command) == "" {
					continue
				}
				if (event == "PreToolUse" || event == "PermissionRequest") && h.PermissionDecision != "" && strings.ToLower(strings.TrimSpace(h.PermissionDecision)) != "deny" {
					h.PermissionDecision = ""
				}
				if (event == "PreToolUse" || event == "PermissionRequest") && h.Decision != "" && strings.ToLower(strings.TrimSpace(h.Decision)) != "deny" {
					h.Decision = ""
				}
				hooks = append(hooks, h)
			}
			if len(hooks) == 0 {
				continue
			}
			kept = append(kept, hookGroup{Matcher: normalizedMatcher, Hooks: hooks})
		}
		if len(kept) > 0 {
			result[event] = kept
		}
	}
	return result
}

func isSupportedHookEvent(event string) bool {
	switch event {
	case "SessionStart", "UserPromptSubmit", "PreToolUse", "PostToolUse", "PermissionRequest", "Stop":
		return true
	default:
		return false
	}
}

func normalizeMatcher(event, matcher string) (string, bool) {
	matcher = strings.TrimSpace(matcher)
	if matcher == "" {
		return "", true
	}
	if event == "SessionStart" {
		return filterMatcher(matcher, []string{"startup", "resume"})
	}
	if event == "PreToolUse" || event == "PostToolUse" || event == "PermissionRequest" {
		return filterMatcher(matcher, []string{"Bash"})
	}
	return matcher, true
}

func filterMatcher(matcher string, allowed []string) (string, bool) {
	parts := strings.Split(matcher, "|")
	allow := map[string]bool{}
	for _, a := range allowed {
		allow[strings.ToLower(strings.TrimSpace(a))] = true
	}
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		token := strings.TrimSpace(p)
		if token == "" {
			continue
		}
		if allow[strings.ToLower(token)] {
			out = append(out, token)
		}
	}
	if len(out) == 0 {
		return "", false
	}
	return strings.Join(out, "|"), true
}

func buildHookWrapperName(event string, groupIdx, hookIdx int, command string) string {
	h := sha1.Sum([]byte(fmt.Sprintf("%s|%d|%d|%s", event, groupIdx, hookIdx, command)))
	hash := hex.EncodeToString(h[:])[:8]
	base := "hook"
	re := regexp.MustCompile(`([a-zA-Z0-9._-]+\.cjs)`)
	matches := re.FindAllString(command, -1)
	if len(matches) > 0 {
		base = strings.TrimSuffix(matches[len(matches)-1], ".cjs")
	}
	return fmt.Sprintf("%s-%s-wrapper.cjs", hash, base)
}

func rewriteHookCommandPathToCodex(command string) string {
	replacements := []struct{ old, new string }{
		{"$HOME/.claude/hooks/", "$HOME/.codex/hooks/"},
		{"~/.claude/hooks/", "~/.codex/hooks/"},
		{".claude/hooks/", ".codex/hooks/"},
		{".claude\\hooks\\", ".codex\\hooks\\"},
	}
	out := command
	for _, r := range replacements {
		out = strings.ReplaceAll(out, r.old, r.new)
	}
	return out
}

func writeHookWrapper(path, originalCommand, event string, overwrite bool) error {
	if _, err := os.Stat(path); err == nil && !overwrite {
		return nil
	}
	js := fmt.Sprintf(`#!/usr/bin/env node
const { spawn } = require("child_process");
const EVENT = %q;
const ORIGINAL_COMMAND = %q;

function sanitize(line) {
  try {
    const obj = JSON.parse(line);
    if (EVENT === "PreToolUse" || EVENT === "PermissionRequest" || EVENT === "Stop") {
      delete obj.additionalContext;
    }
    if (EVENT === "PreToolUse" || EVENT === "PermissionRequest") {
      if (obj.permissionDecision && obj.permissionDecision !== "deny") delete obj.permissionDecision;
      if (obj.decision && obj.decision !== "deny") delete obj.decision;
    }
    return JSON.stringify(obj);
  } catch (e) {
    return line;
  }
}

const child = spawn(ORIGINAL_COMMAND, { shell: true, stdio: ["inherit", "pipe", "inherit"] });
let buf = "";
child.stdout.on("data", (chunk) => {
  buf += chunk.toString();
  let idx = buf.indexOf("\n");
  while (idx >= 0) {
    const line = buf.slice(0, idx);
    buf = buf.slice(idx + 1);
    process.stdout.write(sanitize(line) + "\n");
    idx = buf.indexOf("\n");
  }
});
child.on("close", (code) => {
  if (buf.length > 0) process.stdout.write(sanitize(buf) + "\n");
  process.exit(code || 0);
});
`, event, originalCommand)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(js), 0o755)
}
