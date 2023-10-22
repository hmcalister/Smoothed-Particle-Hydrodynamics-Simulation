package gui

import (
	"hmcalister/SmoothedParticleHydrodynamicsSimulation/config"
	"hmcalister/SmoothedParticleHydrodynamicsSimulation/particle"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

type GUIConfig struct {
	simulationConfig *config.SimulationConfig
	window           *sdl.Window
	surface          *sdl.Surface
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
		simulationConfig.SimulationWidth+2*simulationConfig.SimulationPadding, simulationConfig.SimulationHeight+2*simulationConfig.SimulationPadding,
		sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	guiConfig.surface, err = guiConfig.window.GetSurface()
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
