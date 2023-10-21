package config

import "github.com/creasty/defaults"

type SimulationConfig struct {
	// Simulation Config --------------------------------------------------------------------------

	NumParticles       int     `default:"1000"`
	ParticleMass       float64 `default:"1.0"`
	FluidTargetDensity float64 `default:"1.0"`
	GravityStrength    float64 `default:"1.0"`

	// Smoothing Kernel Config --------------------------------------------------------------------

	SmoothingKernelRadius      float64 `default:"1.0"`
	PressureKernelExponent     int     `default:"2"`
	NearPressureKernelExponent int     `default:"4"`

	// GUI Config ---------------------------------------------------------------------------------

	SimulationWidth  int `default:"512"`
	SimulationHeight int `default:"512"`
	// The padding to add around the simulation, for nice visualizations .
	// Note this is added to the simulation width and height, such that
	// the actual window size is width+padding etc.
	SimulationPadding int `default:"10"`
	FramesPerSecond   int `default:"60"`

	// Spatial Hashing Config ---------------------------------------------------------------------

	// Number of bins to hash cells into.
	// If set to -1 this is set to a number of bins equal to
	// the number of cells that cover the screen
	// (1 + ScreenWidth//SmoothingKernelRadius) * (1 + ScreenHeight//SmoothingKernelRadius)
	SpatialHashingBins int `default:"-1"`
}

func CreateDefaultConfig() *SimulationConfig {
	defaultConfig := &SimulationConfig{}
	defaults.Set(defaultConfig)
	return defaultConfig
}
