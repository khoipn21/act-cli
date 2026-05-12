package kit

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func copyFileIfExists(src, dst string, overwrite bool) error {
	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("required source file not found: %s", src)
	}
	if _, err := os.Stat(dst); err == nil && !overwrite {
		return nil
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
