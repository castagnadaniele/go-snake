package snake

import "math/rand"

type Food struct {
	width  int
	height int
}

func NewFood(width, height int) *Food {
	return &Food{width, height}
}

func (f *Food) Generate(c []Coordinate) Coordinate {
	w := rand.Intn(f.width)
	h := rand.Intn(f.height)
	return Coordinate{w, h}
}
