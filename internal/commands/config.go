package commands

import (
	"fmt"
	"io"
	"sort"

	"act-cli/internal/config"
)

func RunConfig(args []string, w io.Writer) error {
	store, err := config.NewStore()
	if err != nil {
		return err
	}
	cfg, err := store.Load()
	if err != nil {
		return err
	}
	if len(args) == 0 || args[0] == "list" {
		keys := make([]string, 0, len(cfg))
		for k := range cfg {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			fmt.Fprintf(w, "%s=%s\n", k, cfg[k])
		}
		return nil
	}
	switch args[0] {
	case "get":
		if len(args) != 2 {
			return fmt.Errorf("usage: act config get <key>")
		}
		fmt.Fprintln(w, cfg[args[1]])
		return nil
	case "set":
		if len(args) != 3 {
			return fmt.Errorf("usage: act config set <key> <value>")
		}
		cfg[args[1]] = args[2]
		if err := store.Save(cfg); err != nil {
			return err
		}
		fmt.Fprintf(w, "saved %s\n", args[1])
		return nil
	default:
		return fmt.Errorf("unknown config action %q", args[0])
	}
}
