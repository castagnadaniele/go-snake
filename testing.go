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
