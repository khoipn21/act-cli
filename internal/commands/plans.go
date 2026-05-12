package commands

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func RunPlans(args []string, w io.Writer) error {
	if len(args) == 0 || args[0] == "validate" {
		planPath := "plan.md"
		if len(args) >= 2 {
			planPath = args[1]
		}
		return validatePlan(planPath, w)
	}
	return fmt.Errorf("unsupported plans action %q", args[0])
}

func validatePlan(planPath string, w io.Writer) error {
	f, err := os.Open(planPath)
	if err != nil {
		return err
	}
	defer f.Close()

	dir := filepath.Dir(planPath)
	scanner := bufio.NewScanner(f)
	missing := []string{}
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, "./phase-") {
			continue
		}
		start := strings.Index(line, "(./")
		end := strings.Index(line, ")")
		if start == -1 || end == -1 || end <= start+2 {
			continue
		}
		rel := line[start+2 : end]
		path := filepath.Join(dir, rel)
		if _, err := os.Stat(path); err != nil {
			missing = append(missing, rel)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	if len(missing) > 0 {
		return fmt.Errorf("missing phase files: %s", strings.Join(missing, ", "))
	}
	fmt.Fprintf(w, "Plan validation passed: %s\n", planPath)
	return nil
}
