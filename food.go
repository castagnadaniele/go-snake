package snake

import (
	"math/rand"
	"time"
)

const (
	ErrBoardFull = FoodError("snake: food: board full, can not generate food coordinate")
)

// FoodGenerator interface describes a food producer.
//
// Generate should return the coordinate of the next food
// to spawn on the board (this new coordinate should not be contained
// in c), or should return error if the board is full.
type FoodGenerator interface {
	// Generate should return the coordinate of the next food
	// to spawn on the board (this new coordinate should not be contained
	// in c), or should return error if the board is full.
	Generate(c []Coordinate) (Coordinate, error)
}

// Food struct which implements snake food coordinate random generation.
type Food struct {
	width  int
	height int
}

// NewFood returns a pointer to Food and seeds the generator with current time
func NewFood(width, height int) *Food {
	rand.Seed(time.Now().UnixNano())
	return &Food{width, height}
}

// Generate returns a random coordinate for the food which is not in c Coordinates.
// If c length is equal to all the available cells in the board it returns ErrBoardFull.
func (f *Food) Generate(c []Coordinate) (Coordinate, error) {
	if len(c) == f.width*f.height {
		return Coordinate{}, ErrBoardFull
	}
	var foodCoordinate Coordinate
	for ok := true; ok; ok = contains(c, foodCoordinate) {
		w := rand.Intn(f.width)
		h := rand.Intn(f.height)
		foodCoordinate = Coordinate{w, h}
	}
	return foodCoordinate, nil
}

// FoodError type defines food errors
type FoodError string

func (e FoodError) Error() string {
	return string(e)
}
