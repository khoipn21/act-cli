package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"act-cli/internal/kit"
)

func RunAgents(args []string, w io.Writer) error {
	action := "list"
	if len(args) > 0 {
		action = args[0]
	}
	if action != "list" {
		return fmt.Errorf("unsupported agents action %q", action)
	}
	kp, err := kit.ResolveKitPath("")
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(filepath.Join(kp, "claude", "agents"))
	if err != nil {
		return err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}
		names = append(names, strings.TrimSuffix(entry.Name(), ".md"))
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Fprintln(w, name)
	}
	return nil
}
