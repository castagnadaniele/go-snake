package snake

import (
	"fmt"
	"math"
)

const (
	ErrHeadOutOfBoard             = SnakeErr("snake head out of board")
	ErrSnakeMustMoveBeforeGrowing = SnakeErr("snake must move before growing")
)

const (
	upDown    = Up | Down
	leftRight = Left | Right
)

type SnakeErr string

func (e SnakeErr) Error() string {
	return string(e)
}

type SnakeInvalidMoveErr struct {
	Face Direction
	Move Direction
}

func NewSnakeInvalidMoveErr(face, move Direction) error {
	return SnakeInvalidMoveErr{face, move}
}

func (e SnakeInvalidMoveErr) Error() string {
	return fmt.Sprintf("snake can not move %v if facing %v", e.Move, e.Face)
}

type Snake interface {
	// GetCoordinates returns the snake internal coordinates.
	GetCoordinates() []Coordinate

	// Move moves the snake head towards direction d, cutting tail coordinate
	// and appending new coordinate on head. Returns ErrHeadOutOfBoard error
	// when head would move out of the board.
	Move(d Direction) error

	// Grow grows snake tail appending the last cutted tail.
	// If snake did not move before growing, it returns
	// ErrSnakeMustMoveBeforeGrowing error
	Grow() error

	// Face returns where the snake head is facing
	Face() Direction
}

type DefaultSnake struct {
	width         int
	height        int
	coordinates   []Coordinate
	lastTail      *Coordinate
	faceDirection Direction
}

// NewSnake returns a new Snake struct pointer initializing snake coordinates
// and setting width and height of the board.
func NewSnake(width, height int) Snake {
	s := &DefaultSnake{width: width, height: height}
	s.coordinates = make([]Coordinate, 3)
	startX := int(math.Round(float64(s.width) * 0.6))
	startY := int(math.Round(float64(s.height) * 0.5))
	s.coordinates[0] = Coordinate{startX, startY}
	s.coordinates[1] = Coordinate{startX + 1, startY}
	s.coordinates[2] = Coordinate{startX + 2, startY}
	s.faceDirection = Left
	return s
}

func (s *DefaultSnake) GetCoordinates() []Coordinate {
	return s.coordinates
}

func (s *DefaultSnake) Move(d Direction) error {
	head := s.setHead(d)
	if head.X < 0 || head.X >= s.width || head.Y < 0 || head.Y >= s.height {
		return ErrHeadOutOfBoard
	}
	if s.isInvalidMove(d) {
		return NewSnakeInvalidMoveErr(s.faceDirection, d)
	}
	s.lastTail = &s.coordinates[len(s.coordinates)-1]
	s.faceDirection = d
	s.coordinates = append([]Coordinate{{head.X, head.Y}},
		s.coordinates[:len(s.coordinates)-1]...)
	return nil
}

func (s *DefaultSnake) setHead(d Direction) Coordinate {
	head := s.coordinates[0]
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
	return head
}

func (s *DefaultSnake) isInvalidMove(d Direction) bool {
	if d != s.faceDirection &&
		((Has(d, upDown) && Has(s.faceDirection, upDown)) ||
			(Has(d, leftRight) && Has(s.faceDirection, leftRight))) {
		return true
	}
	return false
}

func (s *DefaultSnake) Grow() error {
	if s.lastTail == nil {
		return ErrSnakeMustMoveBeforeGrowing
	}
	s.coordinates = append(s.coordinates, *s.lastTail)
	return nil
}

func (s *DefaultSnake) Face() Direction {
	return s.faceDirection
}
