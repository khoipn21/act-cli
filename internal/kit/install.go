package kit

import (
	"fmt"
	"os"
	"path/filepath"

	"act-cli/internal/fsutil"
)

type InstallOptions struct {
	TargetDir string
	KitPath   string
	Overwrite bool
	Scope     string
}

const (
	ScopeProject = "project"
	ScopeGlobal  = "global"
)

func InstallRuntime(opt InstallOptions) error {
	scope := opt.Scope
	if scope == "" {
		scope = ScopeProject
	}

	src := filepath.Join(opt.KitPath, "claude")
	dst := filepath.Join(opt.TargetDir, ".claude")
	if err := fsutil.CopyTree(src, dst, fsutil.CopyOptions{Overwrite: opt.Overwrite}); err != nil {
		return err
	}

	if scope == ScopeGlobal {
		return nil
	}

	if err := copyFirstExisting([]string{
		filepath.Join(opt.KitPath, "claude", "docs"),
		filepath.Join(opt.KitPath, "docs"),
	}, filepath.Join(dst, "docs"), fsutil.CopyOptions{Overwrite: opt.Overwrite}); err != nil {
		return fmt.Errorf("copy docs: %w", err)
	}
	if err := copyFirstExisting([]string{
		filepath.Join(opt.KitPath, "claude", "plans", "templates"),
		filepath.Join(opt.KitPath, "plans", "templates"),
	}, filepath.Join(dst, "plans", "templates"), fsutil.CopyOptions{Overwrite: opt.Overwrite}); err != nil {
		return fmt.Errorf("copy plan templates: %w", err)
	}
	return nil
}

func copyFirstExisting(candidates []string, dst string, opt fsutil.CopyOptions) error {
	for _, src := range candidates {
		st, err := os.Stat(src)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		if !st.IsDir() {
			continue
		}
		return fsutil.CopyTree(src, dst, opt)
	}
	return nil
}
