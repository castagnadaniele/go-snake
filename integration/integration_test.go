package integration_test

import (
	"testing"
	"time"

	"github.com/castagnadaniele/go-snake"
)

func TestGameTicker(t *testing.T) {
	t.Run("should run game", func(t *testing.T) {
		s := snake.NewSnake(60, 60)
		cloak := snake.NewCloak()
		f := &snake.FoodStub{}
		f.Seed([]snake.FoodStubValue{
			{Coord: snake.Coordinate{X: 0, Y: 0}, Err: nil},
			{Coord: snake.Coordinate{X: 1, Y: 0}, Err: nil},
			{Coord: snake.Coordinate{X: 2, Y: 0}, Err: nil},
			{Coord: snake.Coordinate{X: 3, Y: 0}, Err: nil},
		})
		g := snake.NewGame(s, cloak, f)

		g.Start(2 * time.Millisecond)

		go func() {
			<-time.After(5 * time.Millisecond)
			cloak.Stop()
		}()

		// skip init snake coordinates send
		snake.WaitAndReceiveGameChannels(t, g)
		// skip init food coordinate send
		snake.WaitAndReceiveGameChannels(t, g)

		snake.WaitAndReceiveGameChannels(t, g)
		got, _, _ := snake.WaitAndReceiveGameChannels(t, g)
		want := []snake.Coordinate{
			{X: 34, Y: 30},
			{X: 35, Y: 30},
			{X: 36, Y: 30},
		}
		snake.AssertCoordinates(t, got, want)
	})
}
