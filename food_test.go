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

		c := food.Generate(snakeCoordinates)
		for _, sc := range snakeCoordinates {
			if reflect.DeepEqual(c, sc) {
				t.Errorf("snake coordinates %v should not contain food coordinate %v", snakeCoordinates, c)
			}
		}
	})
}
