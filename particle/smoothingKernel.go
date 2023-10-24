package particle

import (
	"math"
)

const (
	smoothingKernelRadius float64 = 20.0
	exponent              float64 = 2
)

func normalizationCoefficient() float64 {
	return ((exponent*exponent + 3*exponent + 2) / (2 * math.Pi * math.Pow(smoothingKernelRadius, exponent+2)))
}

func smoothingKernel(displacement float64) float64 {
	displacement = min(displacement, smoothingKernelRadius)
	return normalizationCoefficient() * math.Pow(smoothingKernelRadius-displacement, exponent)
}

func smoothingKernelGradientMagnitude(displacement float64) float64 {
	displacement = min(displacement, smoothingKernelRadius)
	return exponent * normalizationCoefficient() * math.Pow(smoothingKernelRadius-displacement, exponent-1)
}
