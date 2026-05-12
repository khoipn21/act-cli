package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"act-cli/internal/version"
)

const releaseCheckTimeout = 1500 * time.Millisecond

type latestReleaseResponse struct {
	TagName string `json:"tag_name"`
}

func maybePrintAutoUpdateNotice(args []string, w io.Writer) {
	if !shouldCheckUpdate(args, w) {
		return
	}

	tag, err := fetchLatestReleaseTag(version.Repo, firstNonEmpty(
		os.Getenv("ACT_GITHUB_TOKEN"),
		os.Getenv("GITHUB_TOKEN"),
		os.Getenv("GH_TOKEN"),
	))
	if err != nil || strings.TrimSpace(tag) == "" {
		return
	}

	current := normalizeSemver(version.Version)
	latest := normalizeSemver(tag)
	cmp, err := compareSemver(current, latest)
	if err != nil || cmp >= 0 {
		return
	}

	renderUpdateNotice(w, version.Version, tag)
}

func shouldCheckUpdate(args []string, w io.Writer) bool {
	_ = w
	if strings.TrimSpace(os.Getenv("ACT_DISABLE_AUTO_UPDATE_CHECK")) != "" {
		return false
	}
	if strings.TrimSpace(os.Getenv("CI")) != "" {
		return false
	}
	if len(args) > 0 && (args[0] == "update" || args[0] == "versions") {
		return false
	}
	return true
}

func fetchLatestReleaseTag(repo, token string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: releaseCheckTimeout}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status %d", res.StatusCode)
	}

	var rel latestReleaseResponse
	if err := json.NewDecoder(res.Body).Decode(&rel); err != nil {
		return "", err
	}
	return strings.TrimSpace(rel.TagName), nil
}

func renderUpdateNotice(w io.Writer, current, latest string) {
	withColor := supportsColor()
	c := func(code, s string) string {
		if !withColor {
			return s
		}
		return "\x1b[" + code + "m" + s + "\x1b[0m"
	}

	fmt.Fprintln(w)
	fmt.Fprintln(w, c("36", "╔════════════════════════════════════════════════════════════╗"))
	fmt.Fprintln(w, c("36", "║ ")+c("93", "New ACT CLI version available")+c("36", "                                 ║"))
	versionLine := fmt.Sprintf("Current: %s  Latest: %s", current, latest)
	padding := 58 - visibleLen(versionLine)
	if padding < 1 {
		padding = 1
	}
	fmt.Fprintf(w, "%s %s%s%s\n", c("36", "║ "), c("37", versionLine), strings.Repeat(" ", padding), c("36", "║"))
	fmt.Fprintln(w, c("36", "║ ")+c("92", "Run `act update` to upgrade now")+c("36", "                                  ║"))
	fmt.Fprintln(w, c("36", "╚════════════════════════════════════════════════════════════╝"))
	fmt.Fprintln(w)
}

func visibleLen(s string) int {
	return len([]rune(s))
}

func supportsColor() bool {
	term := strings.ToLower(strings.TrimSpace(os.Getenv("TERM")))
	return os.Getenv("NO_COLOR") == "" && term != "dumb"
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func normalizeSemver(v string) string {
	v = strings.TrimSpace(v)
	return strings.TrimPrefix(v, "v")
}

func compareSemver(a, b string) (int, error) {
	ap, err := parseSemver(a)
	if err != nil {
		return 0, err
	}
	bp, err := parseSemver(b)
	if err != nil {
		return 0, err
	}
	for i := 0; i < 3; i++ {
		if ap[i] < bp[i] {
			return -1, nil
		}
		if ap[i] > bp[i] {
			return 1, nil
		}
	}
	return 0, nil
}

func parseSemver(v string) ([3]int, error) {
	var out [3]int
	v = strings.SplitN(v, "-", 2)[0]
	parts := strings.Split(v, ".")
	if len(parts) < 2 || len(parts) > 3 {
		return out, fmt.Errorf("invalid semver: %s", v)
	}
	for i := 0; i < 3; i++ {
		if i >= len(parts) {
			out[i] = 0
			continue
		}
		n, err := parsePositiveInt(parts[i])
		if err != nil {
			return out, err
		}
		out[i] = n
	}
	return out, nil
}

func parsePositiveInt(s string) (int, error) {
	var n int
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0, fmt.Errorf("invalid number: %s", s)
		}
		n = n*10 + int(r-'0')
	}
	return n, nil
}
