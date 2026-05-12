package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"act-cli/internal/kit"
)

func RunNew(args []string, w io.Writer) error {
	if len(args) < 1 {
		return errors.New("usage: act new <directory> [--kit <path>] [--force]")
	}
	target := args[0]
	kitPath, overwrite, err := parseKitAndForce(args[1:])
	if err != nil {
		return err
	}
	if err := os.MkdirAll(target, 0o755); err != nil {
		return err
	}
	abs, err := filepath.Abs(target)
	if err != nil {
		return err
	}
	if err := kit.InstallRuntime(kit.InstallOptions{TargetDir: abs, KitPath: kitPath, Overwrite: overwrite}); err != nil {
		return err
	}
	fmt.Fprintf(w, "Created project and installed act-kit at %s\n", abs)
	return nil
}
