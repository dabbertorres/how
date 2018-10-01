package how

import (
	"errors"
	"os"
	"path/filepath"
)

var (
	ErrNotStruct    = errors.New("not a struct")
	ErrInvalidValue = errors.New("invalid value")
)

// GetConfigPath returns a recommended config file location
// $XDG_CONFIG_HOME/<path> if it exists, or $HOME/.config/<path>
// if path is an absolute path, it is returned without modification
func GetConfigPath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	cfgDir := os.Getenv("XDG_CONFIG_HOME")
	if cfgDir == "" {
		cfgDir = filepath.Join(os.Getenv("HOME"), ".config")
	}

	return filepath.Join(cfgDir, path)
}

// CreateConfigFile creates a new file (NOT overwriting if it already exists)
// with config encoded with KeyValueEncoder
func CreateConfigFile(path string, config interface{}) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	return NewKeyValueEncoder(file).Encode(config)
}

