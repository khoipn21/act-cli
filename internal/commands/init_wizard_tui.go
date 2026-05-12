package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"act-cli/internal/kit"
)

func runInitWizard(r io.Reader, w io.Writer, opts initOptions) (initOptions, error) {
	ui := newInitWizardUI(w)
	reader := bufio.NewReader(r)

	ui.header()
	ui.section("Scope")
	scopeInput, err := ui.prompt(reader, "Installation scope [project/global]", opts.scope)
	if err != nil {
		return initOptions{}, err
	}
	if scopeInput != "" {
		scope, scopeErr := normalizeScope(scopeInput)
		if scopeErr != nil {
			return initOptions{}, scopeErr
		}
		opts.scope = scope
	}

	// Keep defaults implicit in interactive flow: no target/source banner.
	if opts.scope == kit.ScopeProject {
		absTarget, absErr := filepath.Abs(opts.target)
		if absErr != nil {
			return initOptions{}, absErr
		}
		if _, statErr := os.Stat(absTarget); os.IsNotExist(statErr) {
			createDir, askErr := ui.askYesNo(reader, "Target directory does not exist. Create it?", false)
			if askErr != nil {
				return initOptions{}, askErr
			}
			if !createDir {
				return initOptions{}, fmt.Errorf("target directory does not exist: %s", absTarget)
			}
			if mkErr := os.MkdirAll(absTarget, 0o755); mkErr != nil {
				return initOptions{}, mkErr
			}
		}
	} else {
		home, homeErr := os.UserHomeDir()
		if homeErr != nil {
			return initOptions{}, homeErr
		}
		opts.target = home
	}

	ui.section("Options")
	overwrite, err := ui.askYesNo(reader, "Overwrite existing files?", opts.overwrite)
	if err != nil {
		return initOptions{}, err
	}
	opts.overwrite = overwrite

	if !opts.installSkillsIsSet {
		installDeps, askErr := ui.askYesNo(reader, "Install skills dependencies now?", false)
		if askErr != nil {
			return initOptions{}, askErr
		}
		opts.installSkillsDeps = installDeps
		opts.installSkillsIsSet = true
		if installDeps {
			// User already confirmed deps install in wizard; skip installer prompt.
			opts.autoApproveDeps = true
		}
	}

	ui.footer()
	opts.interactive = true
	return opts, nil
}

type initWizardUI struct {
	w     io.Writer
	color bool
}

func newInitWizardUI(w io.Writer) initWizardUI {
	term := strings.ToLower(strings.TrimSpace(os.Getenv("TERM")))
	noColor := os.Getenv("NO_COLOR") != ""
	return initWizardUI{w: w, color: !noColor && term != "dumb"}
}

func (ui initWizardUI) header() {
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
	fmt.Fprintf(ui.w, "%s %s\n", ui.paint("│", "36"), ui.paint("ACT Init Wizard • interactive setup for runtime bootstrap", "33"))
	fmt.Fprintln(ui.w, ui.paint("└────────────────────────────────────────────────────────────┘", "36"))
}

func (ui initWizardUI) section(name string) {
	fmt.Fprintf(ui.w, "\n%s\n", ui.paint("["+name+"]", "32"))
}

func (ui initWizardUI) info(msg string) {
	fmt.Fprintf(ui.w, "%s %s\n", ui.paint("::", "36"), msg)
}

func (ui initWizardUI) footer() {
	fmt.Fprintln(ui.w)
	fmt.Fprintln(ui.w, ui.paint("Ready to initialize runtime...", "32"))
}

func (ui initWizardUI) prompt(reader *bufio.Reader, label, defaultValue string) (string, error) {
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

func (ui initWizardUI) askYesNo(reader *bufio.Reader, label string, defaultYes bool) (bool, error) {
	def := "y/N"
	if defaultYes {
		def = "Y/n"
	}
	fmt.Fprintf(ui.w, "%s %s [%s]: ", ui.paint("::", "33"), label, def)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return false, err
	}
	answer := strings.ToLower(strings.TrimSpace(line))
	if answer == "" {
		return defaultYes, nil
	}
	return answer == "y" || answer == "yes", nil
}

func (ui initWizardUI) paint(text, code string) string {
	if !ui.color {
		return text
	}
	return "\x1b[" + code + "m" + text + "\x1b[0m"
}
