package particle

import (
	"hmcalister/SmoothedParticleHydrodynamicsSimulation/config"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
)

type Particle struct {
	Position mat.VecDense
	Velocity mat.VecDense
}

type ParticleCollection struct {
	rng              *rand.Rand
	simulationConfig *config.SimulationConfig
	Particles        []*Particle
}

