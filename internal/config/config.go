package config

import (
	"encoding/json"
	"os"
)

const CONFIGFILENAME = ".gatorconfig.json"

type Config struct {
	Db_url       string `json:"db_url"`
	Current_user string `json:"current_user"`
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return homeDir + "/" + CONFIGFILENAME, nil
}

func Read() (*Config, error) {

	// find config file
	path, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}

	// read from config file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// unmarshal data
	var conf Config

	if err := json.Unmarshal(data, &conf); err != nil {
		return nil, err
	}

	return &conf, nil
}

func (c *Config) SetUser(username string) error {
	c.Current_user = username

	confEncoded, err := json.Marshal(c)
	if err != nil {
		return err
	}

	path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, confEncoded, 0644); err != nil {
		return err
	}

	return nil
}
