package config

import (
	"log"
	"math/rand"
)

type SimulationConfig struct {
	// Simulation Config --------------------------------------------------------------------------

	// Particle properties

	NumParticles int     `default:"1000" yaml:"NumParticles"`
	ParticleMass float64 `default:"1.0" yaml:"ParticleMass"`
	ParticleSize int32   `default:"5" yaml:"ParticleSize"`

	FluidTargetDensity          float64 `default:"1.0" yaml:"FluidTargetDensity"`
	PressureCoefficient         float64 `default:"1.0" yaml:"PressureCoefficient"`
	ViscosityCoefficient        float64 `default:"0.0" yaml:"ViscosityCoefficient"`
	CollisionDampingCoefficient float64 `default:"0.0" yaml:"CollisionDampingCoefficient"`
	GravityStrength             float64 `default:"1.0" yaml:"GravityStrength"`

	// Simulation Meta Config ---------------------------------------------------------------------

	SimulationStepSize         float64 `default:"1.0" yaml:"SimulationStepSize"`
	StepsPerFrame              int     `default:"1" yaml:"StepsPerFrame"`
	SimulationNumWorkerThreads int     `default:"8" yaml:"SimulationNumWorkerThreads"`
	SmoothingKernelRadius      float64 `default:"20" yaml:"SmoothingKernelRadius"`
	// If random seed is set to 0, then a random seed is generated instead
	RandomSeed uint64 `default:"0" yaml:"RandomSeed"`

	// GUI Config ---------------------------------------------------------------------------------

	SimulationWidth  int32   `default:"1024" yaml:"SimulationWidth"`
	SimulationHeight int32   `default:"512" yaml:"SimulationHeight"`
	FramesPerSecond  float64 `default:"60" yaml:"FramesPerSecond"`

	// Spatial Hashing Config ---------------------------------------------------------------------

	// Number of bins to hash cells into.
	// If set to -1 this is set to a number of bins equal to
	// the number of particles
	SpatialHashingBins int `default:"-1" yaml:"SpatialHashingBins"`
}

// Finalize the config by performing any last operations.
//
// e.g. if SpatialHashingBins=-1, replace this with the correct number of bins
func (simulationConfig *SimulationConfig) finalizeConfig() {
	if simulationConfig.SpatialHashingBins == -1 {
		simulationConfig.SpatialHashingBins = 10 * simulationConfig.NumParticles
	}

	if simulationConfig.RandomSeed == 0 {
		simulationConfig.RandomSeed = rand.Uint64()
		log.Printf("RANDOM SEED: %v", simulationConfig.RandomSeed)
	}

}
