package snake

import (
	"math/rand"
	"time"
)

type Food struct {
	width  int
	height int
}

func NewFood(width, height int) *Food {
	rand.Seed(time.Now().UnixNano())
	return &Food{width, height}
}

func (f *Food) Generate(c []Coordinate) Coordinate {
	var foodCoordinate Coordinate
	for ok := true; ok; ok = contains(c, foodCoordinate) {
		w := rand.Intn(f.width)
		h := rand.Intn(f.height)
		foodCoordinate = Coordinate{w, h}
	}
	return foodCoordinate
}
