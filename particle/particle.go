package particle

import "gonum.org/v1/gonum/mat"

type Particle struct {
	Position mat.VecDense
	Velocity mat.VecDense
}
