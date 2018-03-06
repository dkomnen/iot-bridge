package thermometer

import "math/rand"

func randomFloat64InRange(low, high float64) float64 {
	if high <= low {
		return 0
	}
	return (rand.Float64() * low) + (high - low)
}
