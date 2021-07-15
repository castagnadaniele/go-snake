package snake

import (
	"reflect"
	"testing"
)

// AssertCoordinates asserts that got and want coordinates are deep equal
func AssertCoordinates(t testing.TB, got []Coordinate, want []Coordinate) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v coordinates, want %v coordinates", got, want)
	}
}

// AssertNoError asserts that got is nil
func AssertNoError(t testing.TB, got error) {
	t.Helper()
	if got != nil {
		t.Fatalf("got an error but didn't want one: %v", got)
	}
}

// AssertError asserts that got is the error I want
func AssertError(t testing.TB, got error, want error) {
	t.Helper()

	if got == nil {
		t.Fatal("didn't get an error but wanted one")
	}

	if got != want {
		t.Errorf("got error %q, want error %q", got, want)
	}
}

// AssertDirection asserts that the direction I got is the direction I want
func AssertDirection(t testing.TB, got Direction, want Direction) {
	t.Helper()
	if got != want {
		t.Errorf("got %q direction, want %q direction", got, want)
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
