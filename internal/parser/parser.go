package parser

import (
	"encoding/json"
	"os"
)

type Bliss struct {
	Operations []Operation `json:"operations"`
}

type Operation struct {
	Endpoint string `json:"endpoint"`
	Query    string `json:"query"`
	Handler  string `json:"handler"`
}

func GetBliss(path string) (Bliss, error) {
	file, err := os.Open(path)
	if err != nil {
		return Bliss{}, err
	}
	defer file.Close()

	var aha Bliss

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&aha)
	if err != nil {
		return Bliss{}, err
	}
	return aha, nil
}
