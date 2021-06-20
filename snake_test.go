package snake_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/castagnadaniele/go-snake"
)

func TestSnake(t *testing.T) {
	coordinateTests := []struct {
		width            int
		height           int
		snakeCoordinates []snake.Coordinate
	}{
		{60, 60, []snake.Coordinate{{36, 30}, {37, 30}, {38, 30}}},
		{80, 80, []snake.Coordinate{{48, 40}, {49, 40}, {50, 40}}},
		{49, 53, []snake.Coordinate{{29, 27}, {30, 27}, {31, 27}}},
	}

	for _, c := range coordinateTests {
		t.Run(fmt.Sprintf("should start at %v with screen width and size %d, %d", c.snakeCoordinates[0], c.width, c.height), func(t *testing.T) {
			s := snake.NewSnake(c.width, c.height)
			s.Start()

			got := s.Coordinates
			want := c.snakeCoordinates
			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v coordinates, want %v coordinates", got, want)
			}
		})
	}

	t.Run("should have 3 unit length", func(t *testing.T) {
		s := snake.NewSnake(60, 60)
		s.Start()

		got := len(s.Coordinates)
		want := 3
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("should move one cell to left", func(t *testing.T) {
		s := snake.NewSnake(60, 60)
		s.Start()
		s.Move()

		got := s.Coordinates
		want := []snake.Coordinate{
			{35, 30},
			{36, 30},
			{37, 30},
		}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
