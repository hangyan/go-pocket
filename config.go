package main

import (
	"encoding/json"
	"errors"
	"os"
	"os/user"
	"path/filepath"
)

const (
	CONFIG_PATH = ".silencerc"
)

var (
	ErrConfigDoesNotExist = errors.New("config does not exist;try authorize first")
	ErrInvalidConfig      = errors.New("invalid config")
)

type PocketConfig struct {
	ConsumerKey string `json:"consumer_key",omitempty`
	AccessToken string `json:"access_token,omitempty"`
	Username    string `json:"username,omitempty"`
}

func saveConfig(cfg *PocketConfig) error {
	user, err := user.Current()
	if err != nil {
		return err
	}
	path := filepath.Join(user.HomeDir, CONFIG_PATH)
	f, err := os.OpenFile(path, os.O_WRONLY, 0600)
	if err != nil {
		if os.IsNotExist(err) {
			fc, fErr := os.Create(path)
			if fErr != nil {
				return err
			}
			fc.Chmod(0600)
			f = fc
		} else {
			return err
		}
	}
	defer f.Close()

	if err := json.NewEncoder(f).Encode(cfg); err != nil {
		return err
	}
	return nil
}

func loadConfig() (*PocketConfig, error) {
	user, err := user.Current()
	if err != nil {
		return nil, err
	}

	path := filepath.Join(user.HomeDir, CONFIG_PATH)
	f, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, ErrConfigDoesNotExist
		} else {
			return nil, err
		}

	}
	defer f.Close()

	var cfg *PocketConfig
	if err := json.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, ErrInvalidConfig
	}
	return cfg, nil

}
