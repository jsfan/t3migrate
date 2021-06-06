package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

// ReadConfig reads the configuration from a YAML file
func ReadConfig(inputFile string) (*Config, error) {
	inventoryFile, err := os.Open(inputFile)
	if err != nil {
		return nil, fmt.Errorf("failed to open inventory file: %w", err)
	}
	inventoryIn, err := ioutil.ReadAll(inventoryFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read inventory file: %w", err)
	}
	var config *Config
	err = yaml.Unmarshal(inventoryIn, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse inventory file: %w", err)
	}
	return config, nil
}
