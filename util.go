package fake

import (
	"math"
	"math/rand"
	"time"
)

func GenerateRandom(seed int64) *rand.Rand {
	var rnd *rand.Rand
	if !(seed < 0) {
		rnd = rand.New(rand.NewSource(seed))
	} else {
		rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	}

	return rnd
}

func Round(x, unit float64) float64 {
	//https://stackoverflow.com/questions/39544571/golang-round-to-nearest-0-05
	return float64(int64(x/unit+0.5)) * unit
}

func DegToRad(deg float64) float64 {
	return deg * math.Pi / float64(180)
}

func RandBetween(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
