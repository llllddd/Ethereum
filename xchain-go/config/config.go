package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Lock bool `json:"lock"`
}

var Cfg Configuration

// TODO: NEED UNIT TEST
func ParseConfig(configPath string) error {
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&Cfg)
	if err != nil {
		return err
	}
	return nil
}
