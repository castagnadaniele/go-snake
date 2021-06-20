package snake

import "math"

type Snake struct {
	width  int
	height int
}

type Coordinate struct {
	X int
	Y int
}

func NewSnake(width, height int) *Snake {
	return &Snake{width: width, height: height}
}

func (s *Snake) Start() []Coordinate {
	coordinates := make([]Coordinate, 3)
	startX := int(math.Round(float64(s.width) * 0.6))
	startY := int(math.Round(float64(s.height) * 0.5))
	coordinates[0] = Coordinate{startX, startY}
	coordinates[1] = Coordinate{startX + 1, startY}
	coordinates[2] = Coordinate{startX + 2, startY}
	return coordinates
}
