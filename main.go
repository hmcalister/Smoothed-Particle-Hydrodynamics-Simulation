package main

import (
	"flag"
	"log"
	"time"

	gui "hmcalister/SmoothedParticleHydrodynamicsSimulation/GUI"
	"hmcalister/SmoothedParticleHydrodynamicsSimulation/config"
	"hmcalister/SmoothedParticleHydrodynamicsSimulation/particle"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	// The config for the entire simulation
	simulationConfig *config.SimulationConfig

	// The collection of particles
	particleCollection *particle.ParticleCollection

	// The manager for the GUI
	guiConfig *gui.GUIConfig
)

func init() {
	var err error
	configFilePath := flag.String("configFile", "", "Path to the config file. No path results in default config.")
	flag.Parse()

	// Read the config file
	if *configFilePath == "" {
		simulationConfig = config.CreateDefaultConfig()
	} else {
		simulationConfig, err = config.ReadConfigYaml(*configFilePath)
		if err != nil {
			log.Panicf("error during reading config file: %v", err)
		}
	}

	// Create some particles
	particleCollection = particle.CreateParticleCollection(simulationConfig)

	// Start the GUI
	guiConfig, err = gui.InitGUI(simulationConfig)
	if err != nil {
		log.Panicf("error during gui initialization: %v", err)
	}
}

func main() {
	lastFrameTime := time.Unix(0, 0)

GameLoop:
	for {
		// Handle Events
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				break GameLoop
			}
		}

		// Update particle and draw for this frame
		particleCollection.TickParticles()
		guiConfig.DrawParticles(particleCollection.Particles)

		// Handle frame delay for frames per second
		timeToNextFrame := (1 / simulationConfig.FramesPerSecond) - time.Since(lastFrameTime).Seconds()
		if timeToNextFrame > 0 {
			sdl.Delay(uint32(1000 * timeToNextFrame))
		} else {
			log.Printf("FRAME DELAYED")
		}
		lastFrameTime = time.Now()
	}
}
