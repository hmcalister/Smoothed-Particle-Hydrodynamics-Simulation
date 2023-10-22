package main

import (
	"flag"
	"log"

	gui "hmcalister/SmoothedParticleHydrodynamicsSimulation/GUI"
	"hmcalister/SmoothedParticleHydrodynamicsSimulation/config"
)

var (
	// The config for the entire simulation
	simulationConfig *config.SimulationConfig

	// The manager for the GUI
	guiConfig *gui.GUIConfig
)

func init() {
	var err error
	configFilePath := flag.String("ConfigFile", "", "Path to the config file. No path results in default config.")
	flag.Parse()

	if *configFilePath == "" {
		simulationConfig = config.CreateDefaultConfig()
	} else {
		simulationConfig, err = config.ReadConfigYaml(*configFilePath)
		if err != nil {
			log.Panicf("error during reading config file: %v", err)
		}
	}

	guiConfig, err = gui.InitGUI(simulationConfig)
	if err != nil {
		log.Panicf("error during gui initialization: %v", err)
	}
}

func main() {
	// guiConfig.DestroyGUI()
}
