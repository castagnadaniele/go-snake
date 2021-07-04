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
		g := snake.NewGame(s, cloak)

		g.Start(2 * time.Millisecond)

		go func() {
			<-time.After(5 * time.Millisecond)
			cloak.Stop()
		}()

		<-g.Coordinates()
		got := <-g.Coordinates()
		want := []snake.Coordinate{
			{X: 34, Y: 30},
			{X: 35, Y: 30},
			{X: 36, Y: 30},
		}
		snake.AssertCoordinates(t, got, want)
	})
}