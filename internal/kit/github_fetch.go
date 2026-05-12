package kit

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const DefaultGitHubKitRepo = "khoipn21/act-kit"

type PreparedKit struct {
	Path    string
	TempDir string
	Source  string
}

type PrepareOptions struct {
	ExplicitKitPath string
	KitRepo         string
	Out             io.Writer
}

func PrepareKit(opts PrepareOptions) (PreparedKit, error) {
	if localPath, err := ResolveKitPath(opts.ExplicitKitPath); err == nil {
		return PreparedKit{Path: localPath, Source: "local"}, nil
	}

	repo := strings.TrimSpace(opts.KitRepo)
	if repo == "" {
		if env := strings.TrimSpace(os.Getenv("ACT_KIT_REPO")); env != "" {
			repo = env
		} else {
			repo = DefaultGitHubKitRepo
		}
	}

	token := firstNonEmpty(
		os.Getenv("ACT_GITHUB_TOKEN"),
		os.Getenv("GITHUB_TOKEN"),
		os.Getenv("GH_TOKEN"),
	)
	if strings.TrimSpace(token) == "" {
		return PreparedKit{}, fmt.Errorf("local kit source not found and github token is unavailable")
	}

	if err := checkRepoAccess(repo, token); err != nil {
		return PreparedKit{}, err
	}
	if opts.Out != nil {
		fmt.Fprintf(opts.Out, "No local kit found. Fetching %s from GitHub...\n", repo)
	}

	tempDir, err := os.MkdirTemp("", "act-kit-fetch-*")
	if err != nil {
		return PreparedKit{}, err
	}

	zipPath := filepath.Join(tempDir, "kit.zip")
	downloadErr := downloadReleaseZip(repo, token, zipPath)
	if downloadErr != nil {
		if err := downloadZipball(repo, token, zipPath); err != nil {
			_ = os.RemoveAll(tempDir)
			return PreparedKit{}, fmt.Errorf("download kit from github failed: %v; fallback failed: %w", downloadErr, err)
		}
	}

	extractDir := filepath.Join(tempDir, "extracted")
	if err := extractZip(zipPath, extractDir); err != nil {
		_ = os.RemoveAll(tempDir)
		return PreparedKit{}, err
	}
	root, err := findKitRoot(extractDir)
	if err != nil {
		_ = os.RemoveAll(tempDir)
		return PreparedKit{}, err
	}
	return PreparedKit{Path: root, TempDir: tempDir, Source: "github"}, nil
}

func checkRepoAccess(repo, token string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s", repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return nil
	}
	return fmt.Errorf("github token cannot access repo %s (status %d)", repo, res.StatusCode)
}

func downloadReleaseZip(repo, token, dst string) error {
	type asset struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}
	type release struct {
		Assets []asset `json:"assets"`
	}

	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", repo)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("latest release lookup failed (status %d)", res.StatusCode)
	}

	var rel release
	if err := json.NewDecoder(res.Body).Decode(&rel); err != nil {
		return err
	}
	for _, a := range rel.Assets {
		if !strings.HasSuffix(strings.ToLower(a.Name), ".zip") {
			continue
		}
		assetURL := fmt.Sprintf("https://api.github.com/repos/%s/releases/assets/%d", repo, a.ID)
		return downloadURLToFile(assetURL, token, "application/octet-stream", dst)
	}
	return fmt.Errorf("no .zip asset found in latest release")
}

func downloadZipball(repo, token, dst string) error {
	url := fmt.Sprintf("https://api.github.com/repos/%s/zipball", repo)
	return downloadURLToFile(url, token, "application/vnd.github+json", dst)
}

func downloadURLToFile(url, token, accept, dst string) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", accept)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("download failed from %s (status %d)", url, res.StatusCode)
	}
	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, res.Body)
	return err
}

func extractZip(zipPath, dstRoot string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	if err := os.MkdirAll(dstRoot, 0o755); err != nil {
		return err
	}
	base := filepath.Clean(dstRoot) + string(os.PathSeparator)
	for _, f := range r.File {
		target := filepath.Join(dstRoot, f.Name)
		cleanTarget := filepath.Clean(target)
		if !strings.HasPrefix(cleanTarget, base) {
			return fmt.Errorf("invalid zip entry path: %s", f.Name)
		}
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(cleanTarget, 0o755); err != nil {
				return err
			}
			continue
		}
		if err := os.MkdirAll(filepath.Dir(cleanTarget), 0o755); err != nil {
			return err
		}
		in, err := f.Open()
		if err != nil {
			return err
		}
		out, err := os.Create(cleanTarget)
		if err != nil {
			in.Close()
			return err
		}
		if _, err := io.Copy(out, in); err != nil {
			out.Close()
			in.Close()
			return err
		}
		out.Close()
		in.Close()
	}
	return nil
}

func findKitRoot(root string) (string, error) {
	var found string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if _, err := os.Stat(filepath.Join(path, "claude", "metadata.json")); err == nil {
			found = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	if found == "" {
		return "", fmt.Errorf("downloaded kit is invalid: missing claude/metadata.json")
	}
	return found, nil
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}
