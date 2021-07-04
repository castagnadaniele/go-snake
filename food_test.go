package snake_test

import (
	"reflect"
	"testing"

	"github.com/castagnadaniele/go-snake"
)

func TestFood(t *testing.T) {
	t.Run("should generate food coordinate", func(t *testing.T) {
		width, height := 10, 10
		food := snake.NewFood(width, height)
		snakeCoordinates := []snake.Coordinate{
			{0, 0},
			{1, 0},
			{2, 0},
		}

		c, err := food.Generate(snakeCoordinates)
		snake.AssertNoError(t, err)
		for _, sc := range snakeCoordinates {
			if reflect.DeepEqual(c, sc) {
				t.Errorf("snake coordinates %v should not contain food coordinate %v", snakeCoordinates, c)
			}
		}
	})

	t.Run("should generate food where there is no snake coordinate", func(t *testing.T) {
		width, height := 10, 10
		food := snake.NewFood(width, height)
		snakeCoordinates := make([]snake.Coordinate, (width*height)-1)
		index := 0
		for i := 0; i < 10; i++ {
			for j := 0; j < 10; j++ {
				if i == 0 && j == 0 {
					continue
				}
				snakeCoordinates[index] = snake.Coordinate{i, j}
				index++
			}
		}

		c, err := food.Generate(snakeCoordinates)
		snake.AssertNoError(t, err)
		for _, sc := range snakeCoordinates {
			if reflect.DeepEqual(c, sc) {
				t.Errorf("snake coordinates %v should not contain food coordinate %v", snakeCoordinates, c)
			}
		}
	})

	t.Run("should return error when snake coordinates fill the board", func(t *testing.T) {
		width, height := 10, 10
		food := snake.NewFood(width, height)
		snakeCoordinates := make([]snake.Coordinate, (width * height))
		index := 0
		for i := 0; i < 10; i++ {
			for j := 0; j < 10; j++ {
				snakeCoordinates[index] = snake.Coordinate{i, j}
				index++
			}
		}

		_, err := food.Generate(snakeCoordinates)
		snake.AssertError(t, err, snake.ErrBoardFull)
	})
}
