package particle

import (
	"math"
)

const (
	exponent float64 = 2
)

type smoothingKernelStructure struct {
	kernelRadius        float64
	normalizationFactor float64
}

func newSmoothingKernel(kernelRadius float64) *smoothingKernelStructure {
	return &smoothingKernelStructure{
		kernelRadius:        kernelRadius,
		normalizationFactor: ((exponent*exponent + 3*exponent + 2) / (2 * math.Pi * math.Pow(kernelRadius, exponent+2))),
	}
}

func (kernel *smoothingKernelStructure) kernel(displacement float64) float64 {
	displacement = min(displacement, kernel.kernelRadius)
	return kernel.normalizationFactor * math.Pow(kernel.kernelRadius-displacement, exponent)
}

func (kernel *smoothingKernelStructure) kernelGradientMagnitude(displacement float64) float64 {
	displacement = min(displacement, kernel.kernelRadius)
	return exponent * kernel.normalizationFactor * math.Pow(kernel.kernelRadius-displacement, exponent-1)
}
