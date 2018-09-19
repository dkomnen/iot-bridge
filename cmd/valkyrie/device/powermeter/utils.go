package powermeter

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func randomFloat32InRange(low, high float32) float32 {
	return low + rand.Float32()*(high-low)
}
