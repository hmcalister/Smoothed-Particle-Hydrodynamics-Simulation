package particle

import (
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
