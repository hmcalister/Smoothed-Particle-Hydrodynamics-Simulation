package config

import (
	"os"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

func CreateDefaultConfig() *SimulationConfig {
	defaultConfig := &SimulationConfig{}
	defaults.Set(defaultConfig)
	defaultConfig.finalizeConfig()

	return defaultConfig
}

func ReadConfigYaml(yamlFilePath string) (*SimulationConfig, error) {
	// Ensure the file exists and read it

	fileContents, err := os.ReadFile(yamlFilePath)
	if err != nil {
		return nil, err
	}

	simulationConfig := &SimulationConfig{}
	defaults.Set(simulationConfig)

	// Unmarshal the file contents into a SimulationConfig struct
	err = yaml.Unmarshal(fileContents, simulationConfig)
	if err != nil {
		return nil, err
	}

	simulationConfig.finalizeConfig()
	return simulationConfig, nil
}
