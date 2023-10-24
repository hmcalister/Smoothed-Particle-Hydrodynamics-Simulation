package particle

import (
	"gonum.org/v1/gonum/mat"
)

const (
	xDIR int = 0
	yDIR int = 1
)

type Particle struct {
	PredictedPosition *mat.VecDense
	Position          *mat.VecDense
	Velocity          *mat.VecDense
}

func (p *Particle) updatePredictedPosition(stepSize float64) {
	p.PredictedPosition.AddScaledVec(p.Position, stepSize, p.Velocity)
}
