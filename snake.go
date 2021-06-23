package snake

import (
	"math"
)

const (
	ErrHeadOutOfBoard             = SnakeErr("snake head out of board")
	ErrSnakeMustMoveBeforeGrowing = SnakeErr("snake must move before growing")
)

type SnakeErr string

func (e SnakeErr) Error() string {
	return string(e)
}

type Snake struct {
	width       int
	height      int
	coordinates []Coordinate
	lastTail    *Coordinate
}

// NewSnake returns a new Snake struct pointer initializing snake coordinates
// and setting width and height of the board.
func NewSnake(width, height int) *Snake {
	s := &Snake{width: width, height: height}
	s.coordinates = make([]Coordinate, 3)
	startX := int(math.Round(float64(s.width) * 0.6))
	startY := int(math.Round(float64(s.height) * 0.5))
	s.coordinates[0] = Coordinate{startX, startY}
	s.coordinates[1] = Coordinate{startX + 1, startY}
	s.coordinates[2] = Coordinate{startX + 2, startY}
	return s
}

// GetCoordinates returns the snake internal coordinates.
func (s *Snake) GetCoordinates() []Coordinate {
	return s.coordinates
}

// Move moves the snake head towards direction d, cutting tail coordinate
// and appending new coordinate on head. Returns ErrHeadOutOfBoard error
// when head would move out of the board.
func (s *Snake) Move(d Direction) error {
	head := s.coordinates[0]
	s.lastTail = &s.coordinates[len(s.coordinates)-1]
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
	s.coordinates = append([]Coordinate{{head.X, head.Y}},
		s.coordinates[:len(s.coordinates)-1]...)
	return nil
}

// Grow grows snake tail appending the last cutted tail.
// If snake did not move before growing, it returns
// ErrSnakeMustMoveBeforeGrowing error
func (s *Snake) Grow() error {
	if s.lastTail == nil {
		return ErrSnakeMustMoveBeforeGrowing
	}
	s.coordinates = append(s.coordinates, *s.lastTail)
	return nil
}
