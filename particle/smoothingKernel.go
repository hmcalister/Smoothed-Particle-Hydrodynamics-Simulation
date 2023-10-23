package particle

import "math"

func normalizationCoefficient(smoothingRadius float64, exponent float64) float64 {
	return ((exponent*exponent + 3*exponent + 2) / (2 * math.Pi * math.Pow(smoothingRadius, exponent+2)))
}

func smoothingKernel(displacement float64, smoothingRadius float64, exponent float64) float64 {
	return normalizationCoefficient(smoothingRadius, exponent) * math.Pow(displacement, exponent)
}

func smoothingKernelGradientMagnitude(displacement float64, smoothingRadius float64, exponent float64) float64 {
	return exponent * normalizationCoefficient(smoothingRadius, exponent) * math.Pow(displacement, exponent-1)
}
