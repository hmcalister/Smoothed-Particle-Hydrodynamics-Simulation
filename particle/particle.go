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

func CreateInitialParticles(simulationConfig *config.SimulationConfig) *ParticleCollection {
	particleCollection := &ParticleCollection{}
	particleCollection.simulationConfig = simulationConfig
	particleCollection.Particles = make([]*Particle, simulationConfig.NumParticles)

	particleCollection.rng = rand.New(rand.NewSource(simulationConfig.RandomSeed))

	for particleIndex := 0; particleIndex < simulationConfig.NumParticles; particleIndex += 1 {
		particleX := float64(simulationConfig.SimulationWidth) * particleCollection.rng.Float64()
		particleY := float64(simulationConfig.SimulationHeight) * particleCollection.rng.Float64()
		particleCollection.Particles[particleIndex] = &Particle{
			Position: *mat.NewVecDense(2, []float64{particleX, particleY}),
			Velocity: *mat.NewVecDense(2, nil),
		}
	}

	return particleCollection
}
