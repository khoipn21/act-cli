package commands

import (
	"fmt"
	"io"
	"os"
	"strings"
)

type helpTUI struct {
	w     io.Writer
	color bool
}

func renderHelpTUI(w io.Writer, commands []commandHelp) {
	ui := newHelpTUI(w)
	ui.header()
	ui.section("Usage")
	ui.info("act <command> [args]")
	ui.info("help aliases: act help | act --help | act -h")
	ui.section("Commands")
	for _, cmd := range commands {
		ui.command(cmd)
	}
	ui.footer("Tip: run `act <command>` for execution, `act commands` for compact index.")
}

func newHelpTUI(w io.Writer) helpTUI {
	term := strings.ToLower(strings.TrimSpace(os.Getenv("TERM")))
	noColor := os.Getenv("NO_COLOR") != ""
	color := !noColor && term != "dumb"
	return helpTUI{
		w:     w,
		color: color,
	}
}

func (ui helpTUI) header() {
	lines := []string{
		" █████╗  ██████╗████████╗",
		"██╔══██╗██╔════╝╚══██╔══╝",
		"███████║██║        ██║   ",
		"██╔══██║██║        ██║   ",
		"██║  ██║╚██████╗   ██║   ",
		"╚═╝  ╚═╝ ╚═════╝   ╚═╝   ",
	}
	title := "ACT CLI • Commands and Flags Reference"
	fmt.Fprintln(ui.w, ui.paint("┌────────────────────────────────────────────────────────────┐", "36"))
	for i, l := range lines {
		code := "34"
		if i%2 == 1 {
			code = "35"
		}
		fmt.Fprintf(ui.w, "%s %s\n", ui.paint("│", "36"), ui.paint(l, code))
	}
	fmt.Fprintf(ui.w, "%s %s\n", ui.paint("│", "36"), ui.paint(title, "33"))
	fmt.Fprintln(ui.w, ui.paint("└────────────────────────────────────────────────────────────┘", "36"))
}

func (ui helpTUI) section(name string) {
	fmt.Fprintf(ui.w, "\n%s\n", ui.paint("["+name+"]", "32"))
}

func (ui helpTUI) info(line string) {
	fmt.Fprintf(ui.w, "%s %s\n", ui.paint("::", "36"), line)
}

func (ui helpTUI) command(cmd commandHelp) {
	fmt.Fprintf(ui.w, "%s %s\n", ui.paint("•", "35"), ui.paint(cmd.Name, "33"))
	fmt.Fprintf(ui.w, "  %s %s\n", ui.paint("usage:", "36"), cmd.Usage)
	fmt.Fprintf(ui.w, "  %s %s\n", ui.paint("desc :", "36"), cmd.Description)
	if len(cmd.Flags) == 0 {
		return
	}
	fmt.Fprintf(ui.w, "  %s\n", ui.paint("flags:", "36"))
	for _, f := range cmd.Flags {
		fmt.Fprintf(ui.w, "    %-28s %s\n", f.Flag, f.Description)
	}
}

func (ui helpTUI) footer(text string) {
	fmt.Fprintln(ui.w)
	fmt.Fprintln(ui.w, ui.paint(text, "32"))
}

func (ui helpTUI) paint(text, code string) string {
	if !ui.color {
		return text
	}
	return "\x1b[" + code + "m" + text + "\x1b[0m"
}
