package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

type Store struct {
	Path string
}

func NewStore() (*Store, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	return &Store{Path: filepath.Join(home, ".act", "config.json")}, nil
}

func (s *Store) Load() (map[string]string, error) {
	if _, err := os.Stat(s.Path); errors.Is(err, os.ErrNotExist) {
		return map[string]string{}, nil
	}
	bytes, err := os.ReadFile(s.Path)
	if err != nil {
		return nil, err
	}
	if len(bytes) == 0 {
		return map[string]string{}, nil
	}
	var cfg map[string]string
	if err := json.Unmarshal(bytes, &cfg); err != nil {
		return nil, err
	}
	if cfg == nil {
		cfg = map[string]string{}
	}
	return cfg, nil
}

func (s *Store) Save(cfg map[string]string) error {
	if err := os.MkdirAll(filepath.Dir(s.Path), 0o755); err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	bytes = append(bytes, '\n')
	return os.WriteFile(s.Path, bytes, 0o644)
}
