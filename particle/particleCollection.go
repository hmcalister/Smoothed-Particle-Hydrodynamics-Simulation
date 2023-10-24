package particle

import (
	"hmcalister/SmoothedParticleHydrodynamicsSimulation/config"
	"sync"

	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/mat"
)

type ParticleCollection struct {
	rng              *rand.Rand
	simulationConfig *config.SimulationConfig
	spatialHashing   *spatialHashingStructure
	Particles        []*Particle
	densities        []float64
}

func CreateParticleCollection(simulationConfig *config.SimulationConfig) *ParticleCollection {
	particleCollection := &ParticleCollection{}
	particleCollection.simulationConfig = simulationConfig
	particleCollection.Particles = make([]*Particle, simulationConfig.NumParticles)

	particleCollection.rng = rand.New(rand.NewSource(simulationConfig.RandomSeed))
	particleCollection.spatialHashing = createSpatialHashingStructure(
		2*smoothingKernelRadius,
		particleCollection.simulationConfig.SpatialHashingBins,
		simulationConfig.NumParticles,
		simulationConfig.SimulationWidth,
		simulationConfig.SimulationHeight,
	)

	for particleIndex := 0; particleIndex < simulationConfig.NumParticles; particleIndex += 1 {
		particleX := float64(simulationConfig.SimulationWidth) * particleCollection.rng.Float64()
		particleY := float64(simulationConfig.SimulationHeight) * particleCollection.rng.Float64()
		particleCollection.Particles[particleIndex] = &Particle{
			PredictedPosition: mat.NewVecDense(2, []float64{particleX, particleY}),
			Position:          mat.NewVecDense(2, []float64{particleX, particleY}),
			Velocity:          mat.NewVecDense(2, nil),
		}
	}

	particleCollection.densities = make([]float64, simulationConfig.NumParticles)

	return particleCollection
}

func (particleCollection *ParticleCollection) GetParticleColors() []float64 {
	particleColorMap := make([]float64, len(particleCollection.Particles))
	for particleIndex := 0; particleIndex < len(particleCollection.Particles); particleIndex += 1 {
		currentColorMap := (particleCollection.densities[particleIndex] - particleCollection.simulationConfig.FluidTargetDensity) / particleCollection.simulationConfig.FluidTargetDensity
		currentColorMap = min(currentColorMap, 1)
		currentColorMap = max(currentColorMap, -1)
		particleColorMap[particleIndex] = currentColorMap
	}

	// velocityNormValue := 2.0
	// velocityMidValue := 1.0
	// for particleIndex := 0; particleIndex < len(particleCollection.Particles); particleIndex += 1 {
	// 	currentColorMap := (particleCollection.Particles[particleIndex].Velocity.Norm(2) - velocityMidValue) / velocityNormValue
	// 	currentColorMap = min(currentColorMap, 1)
	// 	currentColorMap = max(currentColorMap, -1)
	// 	particleColorMap[particleIndex] = currentColorMap
	// }
	return particleColorMap
}

func (particleCollection *ParticleCollection) calculateDensityWorker(particleIndexChannel <-chan int) {
	displacementVec := mat.NewVecDense(2, nil)
	for particleIndex := range particleIndexChannel {
		density := 0.0

		targetParticle := particleCollection.Particles[particleIndex]
		neighboringParticleIndices := particleCollection.spatialHashing.getAllNeighboringParticleIndices(targetParticle)
		for _, neighborIndex := range neighboringParticleIndices {
			displacementVec.SubVec(targetParticle.PredictedPosition, particleCollection.Particles[neighborIndex].PredictedPosition)
			displacementMagnitude := displacementVec.Norm(2)
			influence := smoothingKernel(displacementMagnitude)
			density += particleCollection.simulationConfig.ParticleMass * influence
		}
		particleCollection.densities[particleIndex] = density
	}
}

func (particleCollection *ParticleCollection) calculateSharedPressure(densityA float64, densityB float64) float64 {
	pressureA := densityA - particleCollection.simulationConfig.FluidTargetDensity
	pressureB := densityB - particleCollection.simulationConfig.FluidTargetDensity
	return particleCollection.simulationConfig.PressureCoefficient * (pressureA + pressureB) / 2
}

func (particleCollection *ParticleCollection) tickParticleWorker(particleIndexChannel <-chan int) {
	totalForce := mat.NewVecDense(2, nil)
	for particleIndex := range particleIndexChannel {
		targetParticle := particleCollection.Particles[particleIndex]

		// Remember - y axis starts with 0 at the top and increases *downwards*
		totalForce.Zero()
		totalForce.SetVec(yDIR, particleCollection.simulationConfig.GravityStrength)

		neighboringParticleIndices := particleCollection.spatialHashing.getAllNeighboringParticleIndices(targetParticle)
		// Calculate influence due to neighboring particles
		displacementVec := mat.NewVecDense(2, nil)
		velocityDifferential := mat.NewVecDense(2, nil)
		for _, neighborIndex := range neighboringParticleIndices {
			if particleIndex == neighborIndex {
				continue
			}

			// Get neighboring particle, find distance to that neighbor
			neighborParticle := particleCollection.Particles[neighborIndex]
			displacementVec.SubVec(targetParticle.PredictedPosition, neighborParticle.PredictedPosition)
			displacementMagnitude := displacementVec.Norm(2)

			// Convert displacement to direction by scaling to unit vector
			displacementVec.ScaleVec(1/displacementMagnitude, displacementVec)

			// Get magnitude of gradient at this displacement
			gradientMagnitude := smoothingKernelGradientMagnitude(displacementMagnitude)

			// Find average pressure between the two particles and use this (approximating newtons third law)
			sharedPressure := particleCollection.calculateSharedPressure(particleCollection.densities[particleIndex], particleCollection.densities[neighborIndex])
			pressureContributionMagnitude := sharedPressure * gradientMagnitude * particleCollection.simulationConfig.ParticleMass / particleCollection.densities[neighborIndex]
			totalForce.AddScaledVec(totalForce, pressureContributionMagnitude, displacementVec)

			// Calculate viscosity force
			velocityDifferential.SubVec(targetParticle.Velocity, neighborParticle.Velocity)
			influence := -smoothingKernel(displacementMagnitude)
			totalForce.AddScaledVec(totalForce, influence*particleCollection.simulationConfig.ViscosityCoefficient, velocityDifferential)
		}
		targetParticle.Velocity.AddScaledVec(targetParticle.Velocity, particleCollection.simulationConfig.SimulationStepSize/particleCollection.densities[particleIndex], totalForce)
		targetParticle.Position.AddScaledVec(targetParticle.Position, particleCollection.simulationConfig.SimulationStepSize, targetParticle.Velocity)

		// Handle edge of simulation
		if targetParticle.Position.AtVec(xDIR) <= 0.0 {
			targetParticle.Position.SetVec(xDIR, particleCollection.rng.Float64())
			targetParticle.Velocity.SetVec(xDIR, -targetParticle.Velocity.AtVec(xDIR)*particleCollection.simulationConfig.CollisionDampingCoefficient)
		} else if targetParticle.Position.AtVec(xDIR) >= float64(particleCollection.simulationConfig.SimulationWidth) {
			targetParticle.Position.SetVec(xDIR, float64(particleCollection.simulationConfig.SimulationWidth)-particleCollection.rng.Float64())
			targetParticle.Velocity.SetVec(xDIR, -targetParticle.Velocity.AtVec(xDIR)*particleCollection.simulationConfig.CollisionDampingCoefficient)
		}
		if targetParticle.Position.AtVec(yDIR) <= 0.0 {
			targetParticle.Position.SetVec(yDIR, particleCollection.rng.Float64())
			targetParticle.Velocity.SetVec(yDIR, -targetParticle.Velocity.AtVec(yDIR)*particleCollection.simulationConfig.CollisionDampingCoefficient)
		} else if targetParticle.Position.AtVec(yDIR) >= float64(particleCollection.simulationConfig.SimulationHeight) {
			targetParticle.Position.SetVec(yDIR, float64(particleCollection.simulationConfig.SimulationHeight)-particleCollection.rng.Float64())
			targetParticle.Velocity.SetVec(yDIR, -targetParticle.Velocity.AtVec(yDIR)*particleCollection.simulationConfig.CollisionDampingCoefficient)
		}
	}
}

func (particleCollection *ParticleCollection) TickParticles() {
	var workerThreadWaitGroup sync.WaitGroup
	var particleIndexChannel chan int

	for _, p := range particleCollection.Particles {
		p.updatePredictedPosition(particleCollection.simulationConfig.SimulationStepSize)
	}

	particleCollection.spatialHashing.updateSpatialHashing(particleCollection.Particles)

	// Recalculate Density Array
	particleIndexChannel = make(chan int, 10)
	for workerThreadIndex := 0; workerThreadIndex < particleCollection.simulationConfig.SimulationNumWorkerThreads; workerThreadIndex++ {
		workerThreadWaitGroup.Add(1)
		go func() {
			particleCollection.calculateDensityWorker(particleIndexChannel)
			workerThreadWaitGroup.Done()
		}()
	}
	for particleIndex := 0; particleIndex < len(particleCollection.Particles); particleIndex += 1 {
		particleIndexChannel <- particleIndex
	}
	close(particleIndexChannel)
	workerThreadWaitGroup.Wait()

	// Tick Particles
	particleIndexChannel = make(chan int, 10)
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
