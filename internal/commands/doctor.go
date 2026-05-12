package commands

import (
	"fmt"
	"io"
	"os"

	"act-cli/internal/kit"
)

func RunDoctor(_ []string, w io.Writer) error {
	fmt.Fprintln(w, "act doctor report")
	if _, err := os.Executable(); err == nil {
		fmt.Fprintln(w, "- executable: OK")
	}
	kp, err := kit.ResolveKitPath("")
	if err != nil {
		fmt.Fprintf(w, "- act-kit: ERROR (%v)\n", err)
		return nil
	}
	fmt.Fprintf(w, "- act-kit: OK (%s)\n", kp)
	return nil
}
