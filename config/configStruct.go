package config

import (
	"math"
	"math/rand"
)

type SimulationConfig struct {
	// Simulation Config --------------------------------------------------------------------------

	NumParticles       int     `default:"1000" yaml:"NumParticles"`
	ParticleMass       float64 `default:"1.0" yaml:"ParticleMass"`
	ParticleSize       int32   `default:"5" yaml:"ParticleSize"`
	FluidTargetDensity float64 `default:"1.0" yaml:"FluidTargetDensity"`
	GravityStrength    float64 `default:"1.0" yaml:"GravityStrength"`
	// If random seed is set to 0, then a random seed is generated instead
	RandomSeed uint64 `default:"0" yaml:"RandomSeed"`

	// Smoothing Kernel Config --------------------------------------------------------------------

	SmoothingKernelRadius      float64 `default:"1.0" yaml:"SmoothingKernelRadius"`
	PressureKernelExponent     int     `default:"2" yaml:"PressureKernelExponent"`
	NearPressureKernelExponent int     `default:"4" yaml:"NearPressureKernelExponent"`

	// GUI Config ---------------------------------------------------------------------------------

	SimulationWidth  int32 `default:"1024" yaml:"SimulationWidth"`
	SimulationHeight int32 `default:"512" yaml:"SimulationHeight"`
	// The padding to add around the simulation for nice visualizations .
	// Note this is added to the simulation width and height such that
	// the actual window size is width+padding etc.
	SimulationPadding int32 `default:"10" yaml:"SimulationPadding"`
	FramesPerSecond   int   `default:"60" yaml:"FramesPerSecond"`

	// Spatial Hashing Config ---------------------------------------------------------------------

	// Number of bins to hash cells into.
	// If set to -1 this is set to a number of bins equal to
	// the number of cells that cover the screen
	// (1 + ScreenWidth//SmoothingKernelRadius) * (1 + ScreenHeight//SmoothingKernelRadius)
	SpatialHashingBins int `default:"-1" yaml:"SpatialHashingBins"`
}

// Finalize the config by performing any last operations.
//
// e.g. if SpatialHashingBins=-1, replace this with the correct number of bins
func (simulationConfig *SimulationConfig) finalizeConfig() {
	if simulationConfig.SpatialHashingBins == -1 {
		simulationCoveringBins := (1 + float64(simulationConfig.SimulationWidth)/simulationConfig.SmoothingKernelRadius) * (1 + float64(simulationConfig.SimulationHeight)/simulationConfig.SmoothingKernelRadius)
		simulationConfig.SpatialHashingBins = int(math.Ceil(simulationCoveringBins))
	}

	if simulationConfig.RandomSeed == 0 {
		simulationConfig.RandomSeed = rand.Uint64()
	}
}
