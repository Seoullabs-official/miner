package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type Config struct {
	Domain      string `json:"domain"`
	TargetMiner string `json:"targetMiner"`
	Validator   string `json:"targetValidator"`
}

func LoadConfig(filename string) (*Config, error) {
	var config Config
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("could not read config file: %w", err)
	}
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}
	return &config, nil
}
