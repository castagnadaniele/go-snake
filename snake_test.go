package snake_test

import (
	"fmt"
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
			snake.AssertCoordinates(t, got, want)
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
		width    int
		height   int
		moves    []snake.Direction
		expected []snake.Coordinate
		err      error
	}{
		{60, 60, []snake.Direction{
			snake.Up,
			snake.Right,
		}, []snake.Coordinate{
			{37, 29},
			{36, 29},
			{36, 30},
		}, nil},
		{60, 60, []snake.Direction{
			snake.Down,
			snake.Left,
		}, []snake.Coordinate{
			{35, 31},
			{36, 31},
			{36, 30},
		}, nil},
		{20, 20, []snake.Direction{
			snake.Up,
			snake.Right,
			snake.Right,
			snake.Right,
			snake.Right,
			snake.Right,
			snake.Right,
			snake.Right,
			snake.Right,
		}, []snake.Coordinate{
			{19, 9},
			{18, 9},
			{17, 9},
		}, snake.ErrHeadOutOfBoard},
		{10, 10, []snake.Direction{
			snake.Left,
			snake.Left,
			snake.Left,
			snake.Left,
			snake.Left,
			snake.Left,
			snake.Left,
		}, []snake.Coordinate{
			{0, 5},
			{1, 5},
			{2, 5},
		}, snake.ErrHeadOutOfBoard},
	}

	for _, m := range moveTests {
		t.Run(fmt.Sprintf("should move %v", m.moves), func(t *testing.T) {
			s := snake.NewSnake(m.width, m.height)
			s.Start()

			for _, move := range m.moves {
				err := s.Move(move)
				if err != nil && err != m.err {
					t.Errorf("got error %v, want error %v", err, m.err)
				}
			}

			got := s.Coordinates
			want := m.expected
			snake.AssertCoordinates(t, got, want)
		})
	}

	t.Run("should grow one cell on the tail", func(t *testing.T) {
		s := snake.NewSnake(60, 60)
		s.Start()
		err := s.Move(snake.Left)
		if err != nil {
			t.Fatalf("got an error but didn't want one, %v", err)
		}
		s.Grow()

		got := s.Coordinates
		want := []snake.Coordinate{{35, 30}, {36, 30}, {37, 30}, {38, 30}}
		snake.AssertCoordinates(t, got, want)
	})
}
