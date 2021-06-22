package snake

import (
	"reflect"
	"testing"
)

func AssertCoordinates(t testing.TB, got []Coordinate, want []Coordinate) {
	t.Helper()

	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v coordinates, want %v coordinates", got, want)
	}
}
