package kit

import (
	"fmt"
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

	if err := fsutil.CopyTree(filepath.Join(opt.KitPath, "docs"), filepath.Join(opt.TargetDir, "docs"), fsutil.CopyOptions{Overwrite: opt.Overwrite}); err != nil {
		return fmt.Errorf("copy docs: %w", err)
	}
	if err := fsutil.CopyTree(filepath.Join(opt.KitPath, "plans", "templates"), filepath.Join(opt.TargetDir, "plans", "templates"), fsutil.CopyOptions{Overwrite: opt.Overwrite}); err != nil {
		return fmt.Errorf("copy plan templates: %w", err)
	}
	return nil
}
