package snake_test

import (
	"testing"
	"time"

	"github.com/castagnadaniele/go-snake"
)

func TestGame(t *testing.T) {
	width, height := 60, 60
	f := snake.NewFood(width, height)

	t.Run("should start game", func(t *testing.T) {
		s := snake.NewSnake(width, height)
		cloak := NewStubCloak()
		defer cloak.Stop()
		g, err := snake.NewGame(s, cloak, f)
		snake.AssertNoError(t, err)
		g.Start(time.Microsecond)

		cloak.AddTick()
		g.ReceiveResult()
		cloak.AddTick()
		got, _ := g.ReceiveResult()
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
		g, err := snake.NewGame(s, cloak, f)
		snake.AssertNoError(t, err)
		g.Start(time.Microsecond)

		cloak.AddTick()
		got, _ := g.ReceiveResult()
		want := []snake.Coordinate{
			{35, 30},
			{36, 30},
			{37, 30},
		}
		snake.AssertCoordinates(t, got, want)

		g.SendMove(snake.Up)
		cloak.AddTick()
		got, _ = g.ReceiveResult()
		want = []snake.Coordinate{
			{35, 29},
			{35, 30},
			{36, 30},
		}
		snake.AssertCoordinates(t, got, want)
	})

	t.Run("snake should keep moving on face direction when an invalid move is sent", func(t *testing.T) {
		s := snake.NewSnake(width, height)
		cloak := NewStubCloak()
		defer cloak.Stop()
		g, err := snake.NewGame(s, cloak, f)
		snake.AssertNoError(t, err)
		g.Start(time.Microsecond)

		g.SendMove(snake.Right)
		cloak.AddTick()
		got, _ := g.ReceiveResult()
		want := []snake.Coordinate{
			{35, 30},
			{36, 30},
			{37, 30},
		}
		snake.AssertCoordinates(t, got, want)
	})

	t.Run("game should end when snake moves out of board", func(t *testing.T) {
		s := snake.NewSnake(10, 10)
		cloak := NewStubCloak()
		defer cloak.Stop()
		g, err := snake.NewGame(s, cloak, f)
		snake.AssertNoError(t, err)
		g.Start(time.Microsecond)

		for i := 0; i < 6; i++ {
			cloak.AddTick()
			g.ReceiveResult()
		}

		cloak.AddTick()
		_, result := g.ReceiveResult()

		if result == nil {
			t.Fatal("should have got a boolean pointer, got nil")
		}

		if *result {
			t.Errorf("got game result %t, want %t", *result, false)
		}
	})

	t.Run("snake should grow after eating food", func(t *testing.T) {
		s := snake.NewSnake(10, 10)
		cloak := NewStubCloak()
		defer cloak.Stop()
		sf := &snake.StubFood{}
		sf.Seed([]snake.Coordinate{{5, 5}, {4, 5}, {3, 5}})
		g, err := snake.NewGame(s, cloak, sf)
		snake.AssertNoError(t, err)
		g.Start(time.Microsecond)
		cloak.AddTick()

		c, r := g.ReceiveResult()
		assertNoGameResult(t, r)
		assertSnakeLength(t, c, 4)
	})

	t.Run("should generate food after snake eats", func(t *testing.T) {
		s := snake.NewSnake(10, 10)
		cloak := NewStubCloak()
		defer cloak.Stop()
		sf := &snake.StubFood{}
		sf.Seed([]snake.Coordinate{{5, 5}, {4, 5}, {3, 5}, {2, 5}})
		g, err := snake.NewGame(s, cloak, sf)
		snake.AssertNoError(t, err)
		g.Start(time.Microsecond)
		cloak.AddTick()

		c, r := g.ReceiveResult()
		assertNoGameResult(t, r)
		assertSnakeLength(t, c, 4)
		cloak.AddTick()
		c, r = g.ReceiveResult()
		assertNoGameResult(t, r)
		assertSnakeLength(t, c, 5)
	})
}

func assertNoGameResult(t testing.TB, result *bool) {
	t.Helper()
	if result != nil {
		t.Fatalf("got result %v, want nil", result)
	}
}

func assertSnakeLength(t testing.TB, c []snake.Coordinate, want int) {
	t.Helper()
	got := len(c)
	if got != want {
		t.Errorf("got snake len %d, want %d", got, want)
	}
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
