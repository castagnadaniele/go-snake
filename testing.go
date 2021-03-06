package snake

import (
	"reflect"
	"testing"
	"time"
)

// AssertCoordinate asserts that got coordinate and want coordiante are equal.
func AssertCoordinate(t testing.TB, got Coordinate, want Coordinate) {
	t.Helper()

	if got.X != want.X || got.Y != want.Y {
		t.Errorf("got %v coordinate, want %v coordinate", got, want)
	}
}

// AssertCoordinates asserts that got and want coordinates are deep equal.
func AssertCoordinates(t testing.TB, got []Coordinate, want []Coordinate) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v coordinates, want %v coordinates", got, want)
	}
}

// AssertNoError asserts that got is nil.
func AssertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("got an error but didn't want one: %v", got)
	}
}

// AssertError asserts that got is the error I want.
func AssertError(t testing.TB, got error, want error) {
	t.Helper()

	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got error %q, want error %q", got, want)
	}
}

// AssertDirection asserts that the direction I got is the direction I want.
func AssertDirection(t testing.TB, got Direction, want Direction) {
	t.Helper()
	if got != want {
		t.Errorf("got %q direction, want %q direction", got, want)
	}
}

// WaitAndReceiveGameChannels returns a ([]Coordinate, *bool, *Coordinate) tuple with snake coordinates
// or game result or food coordinate. It waits to receive values from the game exposed receive channels.
func WaitAndReceiveGameChannels(t testing.TB, g *Game) (snakeCoordinate []Coordinate, gameResult *bool, foodCoordinate *Coordinate) {
	t.Helper()
	select {
	case c := <-g.ReceiveSnakeCoordinates():
		return c, nil, nil
	case r := <-g.ReceiveGameResult():
		return nil, &r, nil
	case f := <-g.ReceiveFoodCoordinate():
		return nil, nil, &f
	case <-time.After(time.Millisecond * 5):
		t.Fatal("got nothing from game channels, want snake coordinates or food coordinate or game result")
		return nil, nil, nil
	}
}

// FoodStubValue stores the coordinate and the error
// returned from FoodStub Generate.
type FoodStubValue struct {
	Coord Coordinate
	Err   error
}

// FoodStub stubs a food generator
type FoodStub struct {
	seedValues []FoodStubValue
}

// Generate returns the first food stub value from FoodStub internal
// array, then pops it from the array.
func (s *FoodStub) Generate(c []Coordinate) (Coordinate, error) {
	result := s.seedValues[0]
	s.seedValues = s.seedValues[1:]
	return result.Coord, result.Err
}

// Seed loads the c food values into FoodStub internal array.
func (s *FoodStub) Seed(c []FoodStubValue) {
	s.seedValues = c
}
