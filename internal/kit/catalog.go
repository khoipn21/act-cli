package kit

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Metadata struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func ResolveKitPath(explicit string) (string, error) {
	if explicit != "" {
		return explicit, nil
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
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", errors.New("local kit source not found")
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
