package snake

import (
	"math"
)

const (
	ErrHeadOutOfBoard = SnakeErr("snake: head out of board")
)

type SnakeErr string

func (e SnakeErr) Error() string {
	return string(e)
}

type Snake struct {
	width       int
	height      int
	Coordinates []Coordinate
	lastTail    Coordinate
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
	tail := Coordinate{startX + 2, startY}
	s.Coordinates[2] = tail
	s.lastTail = tail
}

func (s *Snake) Move(d Direction) error {
	head := s.Coordinates[0]
	s.lastTail = s.Coordinates[len(s.Coordinates)-1]
	switch d {
	case Up:
		head.Y--
	case Down:
		head.Y++
	case Left:
		head.X--
	case Right:
		head.X++
	}
	if head.X < 0 || head.X >= s.width || head.Y < 0 || head.Y >= s.height {
		return ErrHeadOutOfBoard
	}
	s.Coordinates = append([]Coordinate{{head.X, head.Y}}, s.Coordinates[:len(s.Coordinates)-1]...)
	return nil
}

func (s *Snake) Grow() {
	s.Coordinates = append(s.Coordinates, s.lastTail)
}
