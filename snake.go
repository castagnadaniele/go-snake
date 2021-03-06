package snake

import (
	"fmt"
	"math"
)

const (
	ErrHeadOutOfBoard             = SnakeErr("snake: head out of board")
	ErrSnakeMustMoveBeforeGrowing = SnakeErr("snake: must move before growing")
	ErrHeadHitBody                = SnakeErr("snake: head hit body")
)

const (
	upDown    = Up | Down
	leftRight = Left | Right
)

// SnakeErr type defines snake errors
type SnakeErr string

func (e SnakeErr) Error() string {
	return string(e)
}

// SnakeInvalidMoveErr implements invalid move direction error
type SnakeInvalidMoveErr struct {
	Face Direction
	Move Direction
}

func NewSnakeInvalidMoveErr(face, move Direction) error {
	return SnakeInvalidMoveErr{face, move}
}

func (e SnakeInvalidMoveErr) Error() string {
	return fmt.Sprintf("snake: can not move %v if facing %v", e.Move, e.Face)
}

// Snake is the struct which implements the snake behaviour.
type Snake struct {
	width         int
	height        int
	initialLength int
	coordinates   []Coordinate
	lastTail      *Coordinate
	faceDirection Direction
}

// NewSnake returns a new Snake struct pointer initializing snake coordinates
// and setting width and height of the board and snake length 3.
func NewSnake(width, height int) *Snake {
	return NewSnakeOfLength(width, height, 3)
}

// NewSnakeOfLength returns a new Snake struct pointer initializing snake coordinates
// and setting width and height of the board and snake length.
func NewSnakeOfLength(width, height, length int) *Snake {
	s := &Snake{width: width, height: height, initialLength: length}
	s.initCoordinates()
	s.faceDirection = Left
	return s
}

func (s *Snake) initCoordinates() {
	s.coordinates = make([]Coordinate, s.initialLength)
	startX := int(math.Floor(float64(s.width) * 0.6))
	startY := int(math.Floor(float64(s.height) * 0.5))
	for i := 0; i < s.initialLength; i++ {
		s.coordinates[i] = Coordinate{startX + i, startY}
	}
}

// GetCoordinates returns the snake internal coordinates.
func (s *Snake) GetCoordinates() []Coordinate {
	return s.coordinates
}

// Move moves the snake head towards direction d, cutting tail coordinate
// and appending new coordinate on head. Returns ErrHeadOutOfBoard error
// when head would move out of the board. Returns SnakeInvalidMoveErr error if
// direction d is inconsistent with face direction. Returns ErrHeadHitBody error
// if the head would move above a body coordinate.
func (s *Snake) Move(d Direction) error {
	head := s.setHead(d)
	if head.X < 0 || head.X >= s.width || head.Y < 0 || head.Y >= s.height {
		return ErrHeadOutOfBoard
	}
	if !s.IsValidMove(d) {
		return NewSnakeInvalidMoveErr(s.faceDirection, d)
	}
	potentialTail := &s.coordinates[len(s.coordinates)-1]
	nextBodyCoordinates := append([]Coordinate{{head.X, head.Y}}, s.coordinates[:len(s.coordinates)-1]...)
	if contains(nextBodyCoordinates[1:], head) {
		return ErrHeadHitBody
	}
	s.lastTail = potentialTail
	s.faceDirection = d
	s.coordinates = nextBodyCoordinates
	return nil
}

func (s *Snake) setHead(d Direction) Coordinate {
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

// IsValidMove tests if direction is valid for next snake move.
func (s *Snake) IsValidMove(d Direction) bool {
	if d != s.faceDirection &&
		((Has(d, upDown) && Has(s.faceDirection, upDown)) ||
			(Has(d, leftRight) && Has(s.faceDirection, leftRight))) {
		return false
	}
	return true
}

// Grow grows snake tail appending the last cutted tail.
// If snake did not move before growing, it returns
// ErrSnakeMustMoveBeforeGrowing error.
func (s *Snake) Grow() error {
	if s.lastTail == nil {
		return ErrSnakeMustMoveBeforeGrowing
	}
	s.coordinates = append(s.coordinates, *s.lastTail)
	return nil
}

// Face returns where the snake head is facing.
func (s *Snake) Face() Direction {
	return s.faceDirection
}

// Reset resets snake internal coordinates, face direction and last tail coordinate pointer.
func (s *Snake) Reset() {
	s.initCoordinates()
	s.faceDirection = Left
	s.lastTail = nil
}
