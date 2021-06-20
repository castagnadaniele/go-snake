package snake

import "math"

type Snake struct {
	width       int
	height      int
	Coordinates []Coordinate
}

func NewSnake(width, height int) *Snake {
	return &Snake{width: width, height: height}
}

func (s *Snake) Start() {
	s.Coordinates = make([]Coordinate, 3)
	startX := int(math.Round(float64(s.width) * 0.6))
	startY := int(math.Round(float64(s.height) * 0.5))
	s.Coordinates[0] = Coordinate{startX, startY}
	s.Coordinates[1] = Coordinate{startX + 1, startY}
	s.Coordinates[2] = Coordinate{startX + 2, startY}
}

func (s *Snake) Move() {
	head := s.Coordinates[0]
	s.Coordinates = append([]Coordinate{{head.X - 1, head.Y}}, s.Coordinates[:len(s.Coordinates)-1]...)
}
