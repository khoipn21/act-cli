package kit

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Metadata struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func ResolveKitPath(explicit string) (string, error) {
	if strings.TrimSpace(explicit) != "" {
		return validateKitPath(explicit)
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	candidates := []string{
		filepath.Join(cwd, "act-kit"),
		filepath.Join(cwd, "..", "act-kit"),
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err != nil {
			continue
		}
		if kitPath, err := validateKitPath(candidate); err == nil {
			return kitPath, nil
		}
	}
	return "", errors.New("local kit source not found")
}

func validateKitPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	if _, err := os.Stat(filepath.Join(abs, "claude", "metadata.json")); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("invalid act-kit path %q: missing claude/metadata.json", path)
		}
		return "", err
	}
	return abs, nil
}

func LoadMetadata(kitPath string) (Metadata, error) {
	data, err := os.ReadFile(filepath.Join(kitPath, "claude", "metadata.json"))
	if err != nil {
		return Metadata{}, err
	}
	var md Metadata
	if err := json.Unmarshal(data, &md); err != nil {
		return Metadata{}, err
	}
	if md.Name == "" {
		md.Name = "act-kit"
	}
	if md.Version == "" {
		md.Version = "unknown"
	}
	return md, nil
}
