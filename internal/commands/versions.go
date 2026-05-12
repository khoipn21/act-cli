package commands

import (
	"fmt"
	"io"

	"act-cli/internal/version"
)

func RunVersions(_ []string, w io.Writer) error {
	fmt.Fprintf(w, "act-cli %s\n", version.Version)
	return nil
}
