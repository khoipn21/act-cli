package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func runMigrateWizardFull(r io.Reader, w io.Writer, opts migrateOptions) (migrateOptions, error) {
	ui := newMigrateWizardUI(w)
	reader := bufio.NewReader(r)

	ui.header()
	ui.section("Target")
	targetInput, err := ui.prompt(reader, "Migrate target [codex]", "codex")
	if err != nil {
		return migrateOptions{}, err
	}
	targetInput = strings.ToLower(strings.TrimSpace(targetInput))
	if targetInput != "codex" {
		return migrateOptions{}, fmt.Errorf("unsupported migration target %q; expected codex", targetInput)
	}
	opts.to = "codex"

	ui.section("Scope")
	scopeInput, err := ui.prompt(reader, "Installation scope [project/global]", opts.scope)
	if err != nil {
		return migrateOptions{}, err
	}
	scope, scopeErr := normalizeScope(scopeInput)
	if scopeErr != nil {
		return migrateOptions{}, scopeErr
	}
	opts.scope = scope
	opts.scopeIsSet = true
	opts.interactive = true
	ui.footer("Migration selections captured.")
	return opts, nil
}

func runMigrateWizardScopeOnly(r io.Reader, w io.Writer, opts migrateOptions) (migrateOptions, error) {
	ui := newMigrateWizardUI(w)
	reader := bufio.NewReader(r)

	ui.header()
	ui.section("Target")
	ui.info("Migrate target: codex")
	ui.section("Scope")
	scopeInput, err := ui.prompt(reader, "Installation scope [project/global]", opts.scope)
	if err != nil {
		return migrateOptions{}, err
	}
	scope, scopeErr := normalizeScope(scopeInput)
	if scopeErr != nil {
		return migrateOptions{}, scopeErr
	}
	opts.scope = scope
	opts.scopeIsSet = true
	opts.interactive = true
	ui.footer("Scope selected.")
	return opts, nil
}

type migrateWizardUI struct {
	w     io.Writer
	color bool
}

func newMigrateWizardUI(w io.Writer) migrateWizardUI {
	term := strings.ToLower(strings.TrimSpace(os.Getenv("TERM")))
	noColor := os.Getenv("NO_COLOR") != ""
	return migrateWizardUI{w: w, color: !noColor && term != "dumb"}
}

func (ui migrateWizardUI) header() {
	lines := []string{
		" █████╗  ██████╗████████╗",
		"██╔══██╗██╔════╝╚══██╔══╝",
		"███████║██║        ██║   ",
		"██╔══██║██║        ██║   ",
		"██║  ██║╚██████╗   ██║   ",
		"╚═╝  ╚═╝ ╚═════╝   ╚═╝   ",
	}
	fmt.Fprintln(ui.w, ui.paint("┌────────────────────────────────────────────────────────────┐", "36"))
	for i, l := range lines {
		code := "34"
		if i%2 == 1 {
			code = "35"
		}
		fmt.Fprintf(ui.w, "%s %s\n", ui.paint("│", "36"), ui.paint(l, code))
	}
	fmt.Fprintf(ui.w, "%s %s\n", ui.paint("│", "36"), ui.paint("ACT Migrate Wizard • Codex migration setup", "33"))
	fmt.Fprintln(ui.w, ui.paint("└────────────────────────────────────────────────────────────┘", "36"))
}

func (ui migrateWizardUI) section(name string) {
	fmt.Fprintf(ui.w, "\n%s\n", ui.paint("["+name+"]", "32"))
}

func (ui migrateWizardUI) info(msg string) {
	fmt.Fprintf(ui.w, "%s %s\n", ui.paint("::", "36"), msg)
}

func (ui migrateWizardUI) footer(msg string) {
	fmt.Fprintln(ui.w)
	fmt.Fprintln(ui.w, ui.paint(msg, "32"))
}

func (ui migrateWizardUI) prompt(reader *bufio.Reader, label, defaultValue string) (string, error) {
	msg := label
	if strings.TrimSpace(defaultValue) != "" {
		msg = fmt.Sprintf("%s [%s]", label, defaultValue)
	}
	fmt.Fprintf(ui.w, "%s %s: ", ui.paint(">>", "36"), msg)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return "", err
	}
	line = strings.TrimSpace(line)
	if line == "" {
		return defaultValue, nil
	}
	return line, nil
}

func (ui migrateWizardUI) paint(text, code string) string {
	if !ui.color {
		return text
	}
	return "\x1b[" + code + "m" + text + "\x1b[0m"
}
