package metric

import (
	"context"
)

// Values represents metric values
type Values map[string]float64

// Generator interface generates metrics
type Generator interface {
	Generate(context.Context) (Values, error)
}
