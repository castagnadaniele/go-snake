package snake_test

import (
	"testing"
	"time"

	"github.com/castagnadaniele/go-snake"
)

func TestGame(t *testing.T) {
	t.Run("should start game", func(t *testing.T) {
		width, height := 60, 60
		s := snake.NewSnake(width, height)
		cloak := NewStubCloak(2)
		g := snake.NewGame(s, cloak)
		g.Start(time.Microsecond)

		<-g.Coordinates()
		got := <-g.Coordinates()
		want := []snake.Coordinate{
			{34, 30},
			{35, 30},
			{36, 30},
		}
		snake.AssertCoordinates(t, got, want)
	})
}

type StubCloak struct {
	Ticks int
	C     chan time.Time
}

func NewStubCloak(ticks int) *StubCloak {
	ticker := make(chan time.Time, ticks)
	return &StubCloak{ticks, ticker}
}

func (c *StubCloak) Start(d time.Duration) {
	defer close(c.C)
	now := time.Now()
	for i := 0; i < c.Ticks; i++ {
		c.C <- now.Add(time.Duration(i) * d)
	}
}

func (c *StubCloak) Tick() <-chan time.Time {
	return c.C
}

func (c *StubCloak) Stop() {}
