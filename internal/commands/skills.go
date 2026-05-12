package commands

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"

	"act-cli/internal/kit"
)

func RunSkills(args []string, w io.Writer) error {
	action := "list"
	if len(args) > 0 {
		action = args[0]
	}
	if action != "list" {
		return fmt.Errorf("unsupported skills action %q", action)
	}
	kp, err := kit.ResolveKitPath("")
	if err != nil {
		return err
	}
	entries, err := os.ReadDir(filepath.Join(kp, "claude", "skills"))
	if err != nil {
		return err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if entry.Name() == "_shared" || entry.Name() == "common" {
			continue
		}
		names = append(names, entry.Name())
	}
	sort.Strings(names)
	for _, name := range names {
		fmt.Fprintln(w, name)
	}
	return nil
}
