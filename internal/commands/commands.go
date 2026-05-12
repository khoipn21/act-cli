package commands

import (
	"fmt"
	"io"
)

func PrintHelp(w io.Writer) {
	renderHelpTUI(w, commandHelpCatalog())
}

func RunCommands(_ []string, w io.Writer) error {
	for _, cmd := range commandHelpCatalog() {
		fmt.Fprintf(w, "%-10s %s\n", cmd.Name, cmd.Description)
	}
	return nil
}
