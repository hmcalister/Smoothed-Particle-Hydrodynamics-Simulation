package gui

import (
	"fmt"
	"hmcalister/SmoothedParticleHydrodynamicsSimulation/config"
	"hmcalister/SmoothedParticleHydrodynamicsSimulation/particle"
	"log"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type GUIConfig struct {
	simulationConfig *config.SimulationConfig
	window           *sdl.Window
	surface          *sdl.Surface
	renderer         *sdl.Renderer
	font             *ttf.Font
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
		simulationConfig.SimulationWidth+simulationConfig.ParticleSize, simulationConfig.SimulationHeight+simulationConfig.ParticleSize,
		sdl.WINDOW_SHOWN)
	if err != nil {
		return nil, err
	}

	guiConfig.surface, err = guiConfig.window.GetSurface()
	if err != nil {
		return nil, err
	}

	guiConfig.renderer, err = guiConfig.window.GetRenderer()
	if err != nil {
		return nil, err
	}

	err = ttf.Init()
	if err != nil {
		return nil, err
	}

	guiConfig.font, err = ttf.OpenFont("assets/Arial.ttf", 18)
	if err != nil {
		return nil, err
	}

	return guiConfig, nil
}

func (guiConfig *GUIConfig) setColorByParticleColorMap(particleColorMap float64) {
	if particleColorMap < 0 {
		lowerColorChannels := uint8(255 * (1 + particleColorMap))
		guiConfig.renderer.SetDrawColor(lowerColorChannels, lowerColorChannels, 255, 0)
	} else {
		lowerColorChannels := uint8(255 * (1 - particleColorMap))
		guiConfig.renderer.SetDrawColor(255, lowerColorChannels, lowerColorChannels, 0)
	}
}

func (guiConfig *GUIConfig) DrawParticles(particles []*particle.Particle, particleColorMap []float64) {
	for particleIndex, particle := range particles {
		guiConfig.setColorByParticleColorMap(particleColorMap[particleIndex])
		rect := sdl.Rect{
			X: int32(particle.Position.AtVec(0)),
			Y: int32(particle.Position.AtVec(1)),
			W: guiConfig.simulationConfig.ParticleSize,
			H: guiConfig.simulationConfig.ParticleSize,
		}
		guiConfig.renderer.FillRect(&rect)
	}
}

func (guiConfig *GUIConfig) DisplayFPSText(currentFPS float64) {
	textColor := sdl.Color{
		R: 255,
		G: 255,
		B: 255,
	}

	fpsTextSurface, _ := guiConfig.font.RenderUTF8Solid(fmt.Sprintf("FPS: %.2f", currentFPS), textColor)
	defer fpsTextSurface.Free()

	textRect := &sdl.Rect{X: 10, Y: 10, W: fpsTextSurface.W, H: fpsTextSurface.H}

	fpsTextTexture, _ := guiConfig.renderer.CreateTextureFromSurface(fpsTextSurface)
	defer fpsTextTexture.Destroy()
	guiConfig.renderer.Copy(fpsTextTexture, nil, textRect)
}

func (guiConfig *GUIConfig) ShowFrame() {
	guiConfig.renderer.SetDrawColor(0, 0, 0, 255)
	guiConfig.renderer.Present()
	guiConfig.renderer.Clear()

}

func (guiConfig *GUIConfig) DestroyGUI() {
	var err error

	guiConfig.surface.Free()
	err = guiConfig.renderer.Destroy()
	if err != nil {
		log.Panicf("error during gui renderer destruction: %v", err)
	}

	err = guiConfig.window.Destroy()
	if err != nil {
		log.Panicf("error during gui window destruction: %v", err)
	}

	guiConfig.font.Close()

	sdl.Quit()
}
