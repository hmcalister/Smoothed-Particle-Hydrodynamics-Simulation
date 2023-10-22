package gui

import (
	"hmcalister/SmoothedParticleHydrodynamicsSimulation/config"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

type GUIConfig struct {
	simulationConfig *config.SimulationConfig
	window           *sdl.Window
}

func InitGUI(simulationConfig *config.SimulationConfig) (*GUIConfig, error) {
	var err error
	guiConfig := &GUIConfig{}
	guiConfig.simulationConfig = simulationConfig

	err = sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		return nil, err
	}

	guiConfig.window, err = sdl.CreateWindow("Smoothed Particle Hydrodynamics", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED,
		simulationConfig.SimulationWidth+simulationConfig.SimulationPadding, simulationConfig.SimulationHeight+simulationConfig.SimulationPadding,
		sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	return guiConfig, nil
}

func (guiConfig *GUIConfig) DestroyGUI() {
	err := guiConfig.window.Destroy()
	if err != nil {
		log.Panicf("error during gui window destruction: %v", err)
	}

	sdl.Quit()
}
