package app

import (
	"fmt"
	"io"

	"act-cli/internal/commands"
)

func Run(args []string, stdout io.Writer, stderr io.Writer) error {
	maybePrintAutoUpdateNotice(args, stdout)

	if len(args) == 0 {
		commands.PrintHelp(stdout)
		return nil
	}

	switch args[0] {
	case "new":
		return commands.RunNew(args[1:], stdout)
	case "init":
		return commands.RunInit(args[1:], stdout)
	case "migrate":
		return commands.RunMigrate(args[1:], stdout)
	case "doctor":
		return commands.RunDoctor(args[1:], stdout)
	case "config":
		return commands.RunConfig(args[1:], stdout)
	case "skills":
		return commands.RunSkills(args[1:], stdout)
	case "agents":
		return commands.RunAgents(args[1:], stdout)
	case "commands":
		return commands.RunCommands(args[1:], stdout)
	case "plans":
		return commands.RunPlans(args[1:], stdout)
	case "versions":
		return commands.RunVersions(args[1:], stdout)
	case "update":
		return commands.RunUpdate(args[1:], stdout)
	case "help", "--help", "-h":
		commands.PrintHelp(stdout)
		return nil
	default:
		return fmt.Errorf("unknown command %q; run `act help`", args[0])
	}
}
