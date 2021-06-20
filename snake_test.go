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
		name := fmt.Sprintf("should start at %v with screen width and size %d, %d", c.snakeCoordinates[0], c.width, c.height)
		t.Run(name, func(t *testing.T) {
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

	moveTests := []struct {
		moves    []snake.Direction
		expected []snake.Coordinate
	}{
		{[]snake.Direction{
			snake.Up,
			snake.Right,
		}, []snake.Coordinate{
			{37, 29},
			{36, 29},
			{36, 30},
		}},
		{[]snake.Direction{
			snake.Down,
			snake.Left,
		}, []snake.Coordinate{
			{35, 31},
			{36, 31},
			{36, 30},
		}},
		{[]snake.Direction{
			snake.Up,
			snake.Right,
			snake.Right,
			snake.Right,
			snake.Right,
			snake.Right,
		}, []snake.Coordinate{
			{41, 29},
			{40, 29},
			{39, 29},
		}},
	}

	for _, m := range moveTests {
		t.Run(fmt.Sprintf("should move %v", m.moves), func(t *testing.T) {
			s := snake.NewSnake(60, 60)
			s.Start()

			for _, move := range m.moves {
				s.Move(move)
			}

			got := s.Coordinates
			want := m.expected
			if !reflect.DeepEqual(got, want) {
				t.Errorf("got %v, want %v", got, want)
			}
		})
	}
}
