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
		simulationConfig.SimulationWidth, simulationConfig.SimulationHeight,
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

func (guiConfig *GUIConfig) DrawParticles(particles []*particle.Particle) {
	guiConfig.surface.FillRect(nil, 0)

	color := sdl.Color{
		R: 255,
		G: 255,
		B: 255,
	}
	pixel := sdl.MapRGB(guiConfig.surface.Format, color.R, color.G, color.B)
	for _, particle := range particles {
		rect := sdl.Rect{
			X: int32(particle.Position.AtVec(0)),
			Y: int32(particle.Position.AtVec(1)),
			W: guiConfig.simulationConfig.ParticleSize,
			H: guiConfig.simulationConfig.ParticleSize,
		}
		guiConfig.surface.FillRect(&rect, pixel)
	}
	guiConfig.window.UpdateSurface()
}

func (guiConfig *GUIConfig) DestroyGUI() {
	err := guiConfig.window.Destroy()
	if err != nil {
		log.Panicf("error during gui window destruction: %v", err)
	}

	sdl.Quit()
}
