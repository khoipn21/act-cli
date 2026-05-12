package commands

type flagHelp struct {
	Flag        string
	Description string
}

type commandHelp struct {
	Name        string
	Usage       string
	Description string
	Flags       []flagHelp
}

func commandHelpCatalog() []commandHelp {
	return []commandHelp{
		{
			Name:        "new",
			Usage:       "act new <dir> [--kit <path>] [--force]",
			Description: "Create a new project and install ACT runtime.",
			Flags: []flagHelp{
				{Flag: "--kit <path>", Description: "Use a specific kit source path."},
				{Flag: "--force", Description: "Overwrite existing files during install."},
			},
		},
		{
			Name:        "init",
			Usage:       "act init [target-dir] [options]",
			Description: "Initialize ACT runtime in project or global scope.",
			Flags: []flagHelp{
				{Flag: "--kit <path>", Description: "Use a specific kit path instead of auto source."},
				{Flag: "--kit-repo <owner/repo>", Description: "GitHub repo for remote kit fetch."},
				{Flag: "--scope project|global", Description: "Set install scope explicitly."},
				{Flag: "--global, -g", Description: "Shortcut for global scope."},
				{Flag: "--project", Description: "Shortcut for project scope."},
				{Flag: "--interactive, -i", Description: "Force interactive wizard mode."},
				{Flag: "--non-interactive", Description: "Disable prompts."},
				{Flag: "--yes, -y", Description: "Non-interactive defaults."},
				{Flag: "--install-skills-deps", Description: "Install skill dependencies."},
				{Flag: "--skip-skills-deps", Description: "Skip skill dependency install."},
				{Flag: "--force", Description: "Overwrite existing files."},
			},
		},
		{
			Name:        "migrate",
			Usage:       "act migrate [target-dir] [options]",
			Description: "Migrate ACT runtime to Codex-compatible structure.",
			Flags: []flagHelp{
				{Flag: "--to codex", Description: "Migration target. Currently only codex is supported."},
				{Flag: "--scope project|global", Description: "Set migration scope explicitly."},
				{Flag: "--global, -g", Description: "Shortcut for global scope."},
				{Flag: "--project", Description: "Shortcut for project scope."},
				{Flag: "--yes, -y", Description: "Non-interactive mode with defaults."},
				{Flag: "--kit <path>", Description: "Kit source when .claude is missing in target."},
				{Flag: "--force", Description: "Overwrite existing migrated files."},
			},
		},
		{
			Name:        "doctor",
			Usage:       "act doctor",
			Description: "Run local health checks for binary and kit resolution.",
		},
		{
			Name:        "config",
			Usage:       "act config [list|get|set]",
			Description: "Manage ACT config key-value settings.",
			Flags: []flagHelp{
				{Flag: "get <key>", Description: "Read value by key."},
				{Flag: "set <key> <value>", Description: "Persist a config key/value."},
				{Flag: "list", Description: "List all config entries (default action)."},
			},
		},
		{
			Name:        "skills",
			Usage:       "act skills [list]",
			Description: "List available skills from the resolved kit.",
		},
		{
			Name:        "agents",
			Usage:       "act agents [list]",
			Description: "List available agents from the resolved kit.",
		},
		{
			Name:        "commands",
			Usage:       "act commands",
			Description: "Show command index with short descriptions.",
		},
		{
			Name:        "plans",
			Usage:       "act plans [validate] [plan.md]",
			Description: "Validate plan references to phase files.",
		},
		{
			Name:        "versions",
			Usage:       "act versions",
			Description: "Show act-cli version.",
		},
		{
			Name:        "update",
			Usage:       "act update [--check]",
			Description: "Check/update act-cli from GitHub releases.",
			Flags: []flagHelp{
				{Flag: "--check", Description: "Only check for newer version."},
			},
		},
	}
}
