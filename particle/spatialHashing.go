package particle

import (
	"math"
)

type spatialHashingStructure struct {
	// Cell Sizes - equal to the smoothing kernel radius
	cellSize float64

	numCellsX int
	numCellsY int

	// Number of bins to hash cells into
	bins int

	// The partial sums of the number of particles seen up to this bin.
	//
	// After calling updateSpatialHashing this array will hold the number
	// particles seen up to this bin. For example, if we had a partial sums array like
	//
	// `[0,0,2,2,2,5,5]`
	//
	// This may correspond to a denseParticleArray of something like
	//
	// `[4,2,3,1,0]`
	//
	// Note that the second bin (index 1, the first bin with items in it) the entry into the
	// partial sums array is 0, followed by 2. This means there are 2 particles in this bin
	// and the particles start in the dense array at index 0. Likewise, for bin index 4
	// we have an entry of 2 followed by 5, so there are 3 particles which start from index 2.
	//
	// See this link for details: https://matthias-research.github.io/pages/tenMinutePhysics/11-hashing.pdf
	partialSums []int

	// The dense array of particle indices.
	//
	// This holds the particle indices for the spatial hashing
	denseParticleArray []int

	// Particle Hashes
	//
	// A working array to hold the particle hashes and to avoid rehashing particles each step.
	particleHashes []int
}

func createSpatialHashingStructure(cellSize float64, spatialHashingBins int, numParticles int, simulationWidth int32, simulationHeight int32) *spatialHashingStructure {
	numCellsX := int(math.Ceil(float64(simulationWidth) / cellSize))
	numCellsY := int(math.Ceil(float64(simulationHeight) / cellSize))
	return &spatialHashingStructure{
		cellSize:           cellSize,
		numCellsX:          numCellsX,
		numCellsY:          numCellsY,
		bins:               spatialHashingBins,
		partialSums:        make([]int, spatialHashingBins+1),
		denseParticleArray: make([]int, numParticles),
		particleHashes:     make([]int, numParticles),
	}
}

func (sh *spatialHashingStructure) hashCoordinate(x int, y int) int {
	hashCoefficients := []int{92837111, 689287499}
	particleHash := hashCoefficients[xDIR]*x + hashCoefficients[yDIR]*y
	if particleHash < 0 {
		particleHash *= -1
	}
	return particleHash % sh.bins
}

func (sh *spatialHashingStructure) convertParticleToCoordinate(particle *Particle) (int, int) {
	hashingVec := particle.PredictedPosition
	return int(hashingVec.AtVec(xDIR) / sh.cellSize), int(hashingVec.AtVec(yDIR) / sh.cellSize)
}

func (sh *spatialHashingStructure) updateSpatialHashing(particles []*Particle) {
	clear(sh.partialSums)

	// Find count of each bin
	for particleIndex := 0; particleIndex < len(particles); particleIndex += 1 {
		sh.particleHashes[particleIndex] = sh.hashCoordinate(sh.convertParticleToCoordinate(particles[particleIndex]))
		sh.partialSums[sh.particleHashes[particleIndex]] += 1
	}

	// Convert spatial lookup to partial sums
	cumulativeSum := 0
	for binIndex := 0; binIndex < len(sh.partialSums); binIndex += 1 {
		sh.partialSums[binIndex] += cumulativeSum
		cumulativeSum = sh.partialSums[binIndex]
	}

	// Fill in dense particle array using indices of cumulative sum
	for particleIndex := 0; particleIndex < len(particles); particleIndex += 1 {
		particleHash := sh.particleHashes[particleIndex]
		denseIndex := sh.partialSums[particleHash]
		sh.partialSums[sh.particleHashes[particleIndex]] -= 1
		sh.denseParticleArray[denseIndex-1] = particleIndex
	}
}

func (sh *spatialHashingStructure) getParticleIndicesInBin(binIndex int) []int {
	// First, find the start index of this bin in the dense array
	denseStartIndex := sh.partialSums[binIndex]
	// Also find the final bin Index
	denseFinalIndex := sh.partialSums[binIndex+1]

	// Create an array to return the particle indices of this bin, fill it, and return
	binParticleIndices := make([]int, denseFinalIndex-denseStartIndex)
	i := 0
	for itemIndex := denseStartIndex; itemIndex < denseFinalIndex; itemIndex, i = itemIndex+1, i+1 {
		binParticleIndices[i] = sh.denseParticleArray[itemIndex]
	}

	return binParticleIndices
}

func (sh *spatialHashingStructure) getAllNeighboringParticleIndices(particle *Particle) []int {
	neighboringParticleIndices := make([]int, 0)

	// Find the cell coordinates of this particle
	centerCellXCoordinate, centerCellYCoordinate := sh.convertParticleToCoordinate(particle)
	// Then check the cells left, right, up, and down.
	// Skip any cells that lie outside the simulation (think boundaries of screen)
	for dx := -1; dx <= 1; dx += 1 {
		if centerCellXCoordinate+dx < 0 || centerCellXCoordinate+dx > sh.numCellsX {
			continue
		}

		for dy := -1; dy <= 1; dy += 1 {
			if centerCellYCoordinate+dy < 0 || centerCellYCoordinate+dy > sh.numCellsY {
				continue
			}

			// TODO: Fix this bit up, it would be better to collect all the indices first then append at once, but this should be okay...?
			neighboringParticleIndices = append(neighboringParticleIndices, sh.getParticleIndicesInBin(sh.hashCoordinate(centerCellXCoordinate+dx, centerCellYCoordinate+dy))...)
		}
	}

	return neighboringParticleIndices
}
