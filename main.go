package main

import (
	"flag"
	"hmcalister/SmoothedParticleHydrodynamicsSimulation/config"
	"log"
)

var (
	// The config for the entire simulation
	SimulationConfigSettings *config.SimulationConfig
)

func init() {
	configFilePath := flag.String("ConfigFile", "", "Path to the config file. No path results in default config.")
	flag.Parse()

	if *configFilePath == "" {
		SimulationConfigSettings = config.CreateDefaultConfig()
	} else {
		var err error
		SimulationConfigSettings, err = config.ReadConfigYaml(*configFilePath)
		if err != nil {
			log.Panicf("error during reading config file: %v", err)
		}
	}
}

func main() {
	log.Printf("%+v", SimulationConfigSettings)
}
