package parser

import (
	"encoding/json"
	"os"
)

func GetConfig(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var config Config

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

type Config struct {
	BlissPath  string `json:"blissPath"`
	ServerPath string `json:"serverPath"`
	WebPath    string `json:"webPath"`
}
