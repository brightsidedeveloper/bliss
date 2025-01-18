package parser

import (
	"encoding/json"
	"os"
)

type AhaJSON struct {
	Operations []Operation `json:"operations"`
}

type QueryParams map[string]string
type Response map[string]string

type Operation struct {
	Name        string      `json:"name"`
	Endpoint    string      `json:"endpoint"`
	Method      string      `json:"method"`
	QueryParams QueryParams `json:"queryParams"`
	Query       string      `json:"query"`
	Handler     string      `json:"handler"`
	ResType     string      `json:"responseType"`
	Res         Response    `json:"response"`
}

func ParseAhaJSON(path string) (AhaJSON, error) {
	file, err := os.Open(path)
	if err != nil {
		return AhaJSON{}, err
	}
	defer file.Close()

	var aha AhaJSON

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&aha)
	if err != nil {
		return AhaJSON{}, err
	}
	return aha, nil
}
