package thermometer

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func randomFloat64InRange(low, high float64) float64 {
	return low + rand.Float64()*(high-low)
}
