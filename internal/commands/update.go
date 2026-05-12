package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"act-cli/internal/version"
)

type updateOptions struct {
	checkOnly bool
}

type githubRelease struct {
	TagName string               `json:"tag_name"`
	Assets  []githubReleaseAsset `json:"assets"`
}

type githubReleaseAsset struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

func RunUpdate(args []string, w io.Writer) error {
	opts, err := parseUpdateArgs(args)
	if err != nil {
		return err
	}

	token := firstNonEmpty(
		os.Getenv("ACT_GITHUB_TOKEN"),
		os.Getenv("GITHUB_TOKEN"),
		os.Getenv("GH_TOKEN"),
	)
	rel, err := fetchLatestRelease(version.Repo, token)
	if err != nil {
		return err
	}

	current := normalizeSemver(version.Version)
	latest := normalizeSemver(rel.TagName)
	cmp, err := compareSemver(current, latest)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "current: %s\n", version.Version)
	fmt.Fprintf(w, "latest: %s\n", rel.TagName)
	if cmp >= 0 {
		fmt.Fprintln(w, "act is already up to date.")
		return nil
	}
	if opts.checkOnly {
		fmt.Fprintln(w, "update available.")
		return nil
	}

	assetName := releaseAssetName()
	asset, err := findReleaseAsset(rel, assetName)
	if err != nil {
		return err
	}

	execPath, err := os.Executable()
	if err != nil {
		return err
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		execPath, _ = filepath.Abs(execPath)
	}

	downloadPath := execPath + ".new"
	_ = os.Remove(downloadPath)
	if err := downloadReleaseAsset(asset, token, downloadPath); err != nil {
		return err
	}
	if runtime.GOOS != "windows" {
		defer os.Remove(downloadPath)
	}

	if runtime.GOOS == "windows" {
		if err := scheduleWindowsBinaryReplace(execPath, downloadPath); err != nil {
			return err
		}
		fmt.Fprintln(w, "update scheduled. close this shell and run `act versions` in a new shell.")
		return nil
	}

	if err := os.Chmod(downloadPath, 0o755); err != nil {
		return err
	}
	if err := os.Rename(downloadPath, execPath); err != nil {
		return err
	}
	fmt.Fprintln(w, "updated act successfully.")
	return nil
}

func parseUpdateArgs(args []string) (updateOptions, error) {
	opts := updateOptions{}
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--check":
			opts.checkOnly = true
		default:
			return updateOptions{}, fmt.Errorf("unknown option %q", args[i])
		}
	}
	return opts, nil
}

func fetchLatestRelease(repo, token string) (githubRelease, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return githubRelease{}, err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return githubRelease{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return githubRelease{}, fmt.Errorf("failed to fetch latest release (status %d)", res.StatusCode)
	}
	var rel githubRelease
	if err := json.NewDecoder(res.Body).Decode(&rel); err != nil {
		return githubRelease{}, err
	}
	return rel, nil
}

func releaseAssetName() string {
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	return fmt.Sprintf("act-%s-%s%s", runtime.GOOS, runtime.GOARCH, ext)
}

func findReleaseAsset(rel githubRelease, want string) (githubReleaseAsset, error) {
	for _, a := range rel.Assets {
		if a.Name == want {
			return a, nil
		}
	}
	return githubReleaseAsset{}, fmt.Errorf("release asset not found: %s", want)
}

func downloadReleaseAsset(asset githubReleaseAsset, token, dstPath string) error {
	req, err := http.NewRequest(http.MethodGet, asset.URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/octet-stream")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("download release asset failed (status %d)", res.StatusCode)
	}

	out, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer out.Close()
	if _, err := io.Copy(out, res.Body); err != nil {
		return err
	}
	return nil
}

func scheduleWindowsBinaryReplace(execPath, newFile string) error {
	psScriptPath := filepath.Join(os.TempDir(), fmt.Sprintf("act-update-%d.ps1", os.Getpid()))
	ps := fmt.Sprintf(`
$ErrorActionPreference = "Stop"
$pidToWait = %d
$dest = %q
$source = %q
for ($i = 0; $i -lt 120; $i++) {
  if (-not (Get-Process -Id $pidToWait -ErrorAction SilentlyContinue)) { break }
  Start-Sleep -Milliseconds 250
}
Start-Sleep -Milliseconds 200
Copy-Item -LiteralPath $source -Destination $dest -Force
Remove-Item -LiteralPath $source -Force -ErrorAction SilentlyContinue
Remove-Item -LiteralPath $PSCommandPath -Force -ErrorAction SilentlyContinue
`, os.Getpid(), execPath, newFile)
	if err := os.WriteFile(psScriptPath, []byte(ps), 0o600); err != nil {
		return err
	}
	cmd := exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", psScriptPath)
	return cmd.Start()
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
