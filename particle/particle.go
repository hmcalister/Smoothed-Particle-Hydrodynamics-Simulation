package particle

import (
	"hmcalister/SmoothedParticleHydrodynamicsSimulation/config"
	"sync"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
)

const (
	xDIR int = 0
	yDIR int = 1
)

type Particle struct {
	Position *mat.VecDense
	Velocity *mat.VecDense
}

type ParticleCollection struct {
	rng              *rand.Rand
	simulationConfig *config.SimulationConfig
	spatialHashing   *spatialHashingStructure
	Particles        []*Particle
}

func CreateParticleCollection(simulationConfig *config.SimulationConfig) *ParticleCollection {
	particleCollection := &ParticleCollection{}
	particleCollection.simulationConfig = simulationConfig
	particleCollection.Particles = make([]*Particle, simulationConfig.NumParticles)

	particleCollection.rng = rand.New(rand.NewSource(simulationConfig.RandomSeed))
	particleCollection.spatialHashing = createSpatialHashingStructure(
		2*simulationConfig.SmoothingKernelRadius,
		particleCollection.simulationConfig.SpatialHashingBins,
		simulationConfig.NumParticles,
		simulationConfig.SimulationWidth,
		simulationConfig.SimulationHeight,
	)

	for particleIndex := 0; particleIndex < simulationConfig.NumParticles; particleIndex += 1 {
		particleX := float64(simulationConfig.SimulationWidth) * particleCollection.rng.Float64()
		particleY := float64(simulationConfig.SimulationHeight) * particleCollection.rng.Float64()
		particleCollection.Particles[particleIndex] = &Particle{
			Position: mat.NewVecDense(2, []float64{particleX, particleY}),
			Velocity: mat.NewVecDense(2, nil),
		}
	}

	return particleCollection
}

func (particleCollection *ParticleCollection) tickParticleWorker(particleIndexChannel <-chan int) {
	for particleIndex := range particleIndexChannel {
		targetParticle := particleCollection.Particles[particleIndex]
		totalAcceleration := mat.NewVecDense(2, nil)

		// neighboringParticleIndices := particleCollection.spatialHashing.getAllNeighboringParticleIndices(targetParticle)
		// log.Printf("PARTICLE INDEX: %v\tNEIGHBORS: %v", particleIndex, neighboringParticleIndices)

		// Remember - y axis starts with 0 at the top and increases *downwards*
		totalAcceleration.SetVec(yDIR, particleCollection.simulationConfig.GravityStrength)

		// TODO: Find All neighboring particles and apply viscosity etc

		targetParticle.Velocity.AddScaledVec(targetParticle.Velocity, particleCollection.simulationConfig.SimulationStepSize, totalAcceleration)
		targetParticle.Position.AddScaledVec(targetParticle.Position, particleCollection.simulationConfig.SimulationStepSize, targetParticle.Velocity)

		// Handle edge of simulation
		if targetParticle.Position.AtVec(xDIR) < 0.0 {
			targetParticle.Position.SetVec(xDIR, -targetParticle.Position.AtVec(xDIR))
			targetParticle.Velocity.SetVec(xDIR, -targetParticle.Velocity.AtVec(xDIR)*particleCollection.simulationConfig.CollisionDampingCoefficient)
		} else if targetParticle.Position.AtVec(xDIR) > float64(particleCollection.simulationConfig.SimulationWidth) {
			targetParticle.Position.SetVec(xDIR, 2*float64(particleCollection.simulationConfig.SimulationWidth)-targetParticle.Position.AtVec(xDIR))
			targetParticle.Velocity.SetVec(xDIR, -targetParticle.Velocity.AtVec(xDIR)*particleCollection.simulationConfig.CollisionDampingCoefficient)
		}
		if targetParticle.Position.AtVec(yDIR) < 0.0 {
			targetParticle.Position.SetVec(yDIR, -targetParticle.Position.AtVec(yDIR))
			targetParticle.Velocity.SetVec(yDIR, -targetParticle.Velocity.AtVec(yDIR)*particleCollection.simulationConfig.CollisionDampingCoefficient)
		} else if targetParticle.Position.AtVec(yDIR) > float64(particleCollection.simulationConfig.SimulationHeight) {
			targetParticle.Position.SetVec(yDIR, 2*float64(particleCollection.simulationConfig.SimulationHeight)-targetParticle.Position.AtVec(yDIR))
			targetParticle.Velocity.SetVec(yDIR, -targetParticle.Velocity.AtVec(yDIR)*particleCollection.simulationConfig.CollisionDampingCoefficient)
		}
	}
}

func (particleCollection *ParticleCollection) TickParticles() {
	var workerThreadWaitGroup sync.WaitGroup

	particleCollection.spatialHashing.updateSpatialHashing(particleCollection.Particles)
	// log.Printf("PARTIAL SUMS:\t%v", particleCollection.spatialHashing.partialSums)
	// log.Printf("DENSE ARRAY:\t%v", particleCollection.spatialHashing.denseParticleArray)

	particleIndexChannel := make(chan int, 10)
	for workerThreadIndex := 0; workerThreadIndex < particleCollection.simulationConfig.SimulationNumWorkerThreads; workerThreadIndex++ {
		workerThreadWaitGroup.Add(1)
		go func() {
			particleCollection.tickParticleWorker(particleIndexChannel)
			workerThreadWaitGroup.Done()
		}()
	}

	for particleIndex := 0; particleIndex < len(particleCollection.Particles); particleIndex += 1 {
		particleIndexChannel <- particleIndex
	}
	close(particleIndexChannel)
	workerThreadWaitGroup.Wait()
}
