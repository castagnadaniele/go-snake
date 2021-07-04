package snake_test

import (
	"testing"
	"time"

	"github.com/castagnadaniele/go-snake"
)

func TestGame(t *testing.T) {
	width, height := 60, 60

	t.Run("should start game", func(t *testing.T) {
		s := snake.NewSnake(width, height)
		cloak := NewStubCloak()
		defer cloak.Stop()
		g := snake.NewGame(s, cloak)
		g.Start(time.Microsecond)

		cloak.AddTick()
		<-g.Coordinates()
		cloak.AddTick()
		got := <-g.Coordinates()
		want := []snake.Coordinate{
			{34, 30},
			{35, 30},
			{36, 30},
		}
		snake.AssertCoordinates(t, got, want)
	})

	t.Run("snake should change direction", func(t *testing.T) {
		s := snake.NewSnake(width, height)
		cloak := NewStubCloak()
		defer cloak.Stop()
		g := snake.NewGame(s, cloak)
		g.Start(time.Microsecond)

		cloak.AddTick()
		got := <-g.Coordinates()
		want := []snake.Coordinate{
			{35, 30},
			{36, 30},
			{37, 30},
		}
		snake.AssertCoordinates(t, got, want)

		g.SendMove(snake.Up)
		cloak.AddTick()
		got = <-g.Coordinates()
		want = []snake.Coordinate{
			{35, 29},
			{35, 30},
			{36, 30},
		}
		snake.AssertCoordinates(t, got, want)
	})
}

type StubCloak struct {
	C        chan time.Time
	now      time.Time
	i        int
	duration time.Duration
}

func NewStubCloak() *StubCloak {
	ticker := make(chan time.Time)
	return &StubCloak{ticker, time.Now(), 0, time.Nanosecond}
}

func (c *StubCloak) Start(d time.Duration) {
	c.duration = d
}

func (c *StubCloak) AddTick() {
	c.C <- c.now.Add(time.Duration(c.i) * c.duration)
}

func (c *StubCloak) Tick() <-chan time.Time {
	return c.C
}

func (c *StubCloak) Stop() {
	close(c.C)
}
