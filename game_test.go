package snake_test

import (
	"testing"
	"time"

	"github.com/castagnadaniele/go-snake"
)

func TestGame(t *testing.T) {
	width, height := 60, 60
	fs := &snake.FoodStub{}
	foodSeededCoordinates := []snake.FoodStubValue{
		{snake.Coordinate{0, 0}, nil},
		{snake.Coordinate{1, 0}, nil},
		{snake.Coordinate{2, 0}, nil},
		{snake.Coordinate{3, 0}, nil},
		{snake.Coordinate{4, 0}, nil},
		{snake.Coordinate{5, 0}, nil},
		{snake.Coordinate{6, 0}, nil},
	}
	fs.Seed(foodSeededCoordinates)

	t.Run("should start game", func(t *testing.T) {
		s := snake.NewSnake(width, height)
		snakeInitCoordinates := s.GetCoordinates()
		cloak := NewStubCloak()
		defer cloak.Stop()
		g := snake.NewGame(s, cloak, fs)
		g.Start(time.Microsecond)

		sc, _, _ := snake.WaitAndReceiveGameChannels(t, g)
		snake.AssertCoordinates(t, sc, snakeInitCoordinates)
		_, _, fc := snake.WaitAndReceiveGameChannels(t, g)
		snake.AssertCoordinate(t, *fc, foodSeededCoordinates[0].Coord)
		cloak.AddTick()
		snake.WaitAndReceiveGameChannels(t, g)
		cloak.AddTick()
		got, r, _ := snake.WaitAndReceiveGameChannels(t, g)
		want := []snake.Coordinate{
			{34, 30},
			{35, 30},
			{36, 30},
		}
		assertNoGameResult(t, r)
		snake.AssertCoordinates(t, got, want)
	})

	t.Run("snake should change direction", func(t *testing.T) {
		s := snake.NewSnake(width, height)
		cloak := NewStubCloak()
		defer cloak.Stop()
		g := snake.NewGame(s, cloak, fs)
		g.Start(time.Microsecond)

		// skip init snake coordinates send
		snake.WaitAndReceiveGameChannels(t, g)
		// skip init food coordinate send
		snake.WaitAndReceiveGameChannels(t, g)

		cloak.AddTick()

		got, r, _ := snake.WaitAndReceiveGameChannels(t, g)
		want := []snake.Coordinate{
			{35, 30},
			{36, 30},
			{37, 30},
		}
		assertNoGameResult(t, r)
		snake.AssertCoordinates(t, got, want)

		g.SendMove(snake.Up)
		cloak.AddTick()

		got, r, _ = snake.WaitAndReceiveGameChannels(t, g)
		want = []snake.Coordinate{
			{35, 29},
			{35, 30},
			{36, 30},
		}
		assertNoGameResult(t, r)
		snake.AssertCoordinates(t, got, want)
	})

	t.Run("snake should keep moving on face direction when an invalid move is sent", func(t *testing.T) {
		s := snake.NewSnake(width, height)
		cloak := NewStubCloak()
		defer cloak.Stop()
		g := snake.NewGame(s, cloak, fs)
		g.Start(time.Microsecond)

		// skip init snake coordinates send
		snake.WaitAndReceiveGameChannels(t, g)
		// skip init food coordinate send
		snake.WaitAndReceiveGameChannels(t, g)

		g.SendMove(snake.Right)
		cloak.AddTick()

		got, r, _ := snake.WaitAndReceiveGameChannels(t, g)
		want := []snake.Coordinate{
			{35, 30},
			{36, 30},
			{37, 30},
		}
		assertNoGameResult(t, r)
		snake.AssertCoordinates(t, got, want)
	})

	t.Run("game should end with a lose when snake moves out of board", func(t *testing.T) {
		s := snake.NewSnake(10, 10)
		cloak := NewStubCloak()
		defer cloak.Stop()
		g := snake.NewGame(s, cloak, fs)
		g.Start(time.Microsecond)

		// skip init snake coordinates send
		snake.WaitAndReceiveGameChannels(t, g)
		// skip init food coordinate send
		snake.WaitAndReceiveGameChannels(t, g)

		for i := 0; i < 6; i++ {
			cloak.AddTick()
			snake.WaitAndReceiveGameChannels(t, g)
		}

		cloak.AddTick()

		_, result, _ := snake.WaitAndReceiveGameChannels(t, g)
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
		sf := &snake.FoodStub{}
		sf.Seed([]snake.FoodStubValue{
			{snake.Coordinate{5, 5}, nil},
			{snake.Coordinate{4, 5}, nil},
			{snake.Coordinate{3, 5}, nil},
		})
		g := snake.NewGame(s, cloak, sf)
		g.Start(time.Microsecond)

		// skip init snake coordinates send
		snake.WaitAndReceiveGameChannels(t, g)
		// skip init food coordinate send
		snake.WaitAndReceiveGameChannels(t, g)

		cloak.AddTick()

		c, r, fc := snake.WaitAndReceiveGameChannels(t, g)
		assertNoGameResult(t, r)
		assertNoFoodCoordinate(t, fc)
		assertSnakeLength(t, c, 4)

		c, r, fc = snake.WaitAndReceiveGameChannels(t, g)
		snake.AssertCoordinate(t, *fc, snake.Coordinate{4, 5})
		assertNoGameResult(t, r)
		assertNoSnakeCoordinates(t, c)
	})

	t.Run("should generate food after snake eats", func(t *testing.T) {
		s := snake.NewSnake(10, 10)
		cloak := NewStubCloak()
		defer cloak.Stop()
		sf := &snake.FoodStub{}
		sf.Seed([]snake.FoodStubValue{
			{snake.Coordinate{5, 5}, nil},
			{snake.Coordinate{4, 5}, nil},
			{snake.Coordinate{3, 5}, nil},
			{snake.Coordinate{2, 5}, nil},
		})
		g := snake.NewGame(s, cloak, sf)
		g.Start(time.Microsecond)

		// skip init snake coordinates send
		snake.WaitAndReceiveGameChannels(t, g)
		// skip init food coordinate send
		snake.WaitAndReceiveGameChannels(t, g)

		cloak.AddTick()

		c, r, fc := snake.WaitAndReceiveGameChannels(t, g)
		assertNoGameResult(t, r)
		assertNoFoodCoordinate(t, fc)
		assertSnakeLength(t, c, 4)

		c, r, fc = snake.WaitAndReceiveGameChannels(t, g)
		assertNoGameResult(t, r)
		assertNoSnakeCoordinates(t, c)
		snake.AssertCoordinate(t, *fc, snake.Coordinate{4, 5})

		cloak.AddTick()

		c, r, fc = snake.WaitAndReceiveGameChannels(t, g)
		assertNoGameResult(t, r)
		assertSnakeLength(t, c, 5)
		assertNoFoodCoordinate(t, fc)

		c, r, fc = snake.WaitAndReceiveGameChannels(t, g)
		assertNoGameResult(t, r)
		assertNoSnakeCoordinates(t, c)
		snake.AssertCoordinate(t, *fc, snake.Coordinate{3, 5})
	})

	t.Run("game should end with a win when snake fills the entire board", func(t *testing.T) {
		s := snake.NewSnakeOfLength(2, 2, 1)
		cloak := NewStubCloak()
		defer cloak.Stop()
		sf := &snake.FoodStub{}
		sf.Seed([]snake.FoodStubValue{
			{snake.Coordinate{0, 1}, nil},
			{snake.Coordinate{0, 0}, nil},
			{snake.Coordinate{1, 0}, nil},
			{snake.Coordinate{1, 1}, snake.ErrBoardFull},
		})
		g := snake.NewGame(s, cloak, sf)
		g.Start(time.Microsecond)

		// skip init snake coordinates send
		snake.WaitAndReceiveGameChannels(t, g)
		// skip init food coordinate send
		snake.WaitAndReceiveGameChannels(t, g)

		cloak.AddTick()

		c, r, fc := snake.WaitAndReceiveGameChannels(t, g)
		assertNoGameResult(t, r)
		assertNoFoodCoordinate(t, fc)
		assertSnakeLength(t, c, 2)

		c, r, fc = snake.WaitAndReceiveGameChannels(t, g)
		assertNoGameResult(t, r)
		assertNoSnakeCoordinates(t, c)
		snake.AssertCoordinate(t, *fc, snake.Coordinate{0, 0})

		g.SendMove(snake.Up)
		cloak.AddTick()

		c, r, fc = snake.WaitAndReceiveGameChannels(t, g)
		assertNoGameResult(t, r)
		assertNoFoodCoordinate(t, fc)
		assertSnakeLength(t, c, 3)

		c, r, fc = snake.WaitAndReceiveGameChannels(t, g)
		assertNoGameResult(t, r)
		assertNoSnakeCoordinates(t, c)
		snake.AssertCoordinate(t, *fc, snake.Coordinate{1, 0})

		g.SendMove(snake.Right)
		cloak.AddTick()

		c, r, fc = snake.WaitAndReceiveGameChannels(t, g)
		assertNoFoodCoordinate(t, fc)
		assertNoGameResult(t, r)
		expectedSnakeCoordinates := []snake.Coordinate{
			{1, 0},
			{0, 0},
			{0, 1},
			{1, 1},
		}
		snake.AssertCoordinates(t, c, expectedSnakeCoordinates)
		assertSnakeLength(t, c, 4)

		c, r, fc = snake.WaitAndReceiveGameChannels(t, g)
		assertNoFoodCoordinate(t, fc)
		assertNoSnakeCoordinates(t, c)
		if r == nil || *r == false {
			t.Errorf("got result %v, want result %t", r, true)
		}
	})

	t.Run("should restart game", func(t *testing.T) {
		s := snake.NewSnake(width, height)
		sf := &snake.FoodStub{}
		sf.Seed(foodSeededCoordinates)
		cloak := NewStubCloak()
		defer cloak.Stop()

		g := snake.NewGame(s, cloak, sf)
		g.Start(time.Microsecond)

		wantSnakeCoord, _, _ := snake.WaitAndReceiveGameChannels(t, g)
		_, _, wantFoodCoord := snake.WaitAndReceiveGameChannels(t, g)

		sf.Seed(foodSeededCoordinates)
		g.Restart(time.Microsecond)

		gotSnakeCoord, _, _ := snake.WaitAndReceiveGameChannels(t, g)
		_, _, gotFoodCoord := snake.WaitAndReceiveGameChannels(t, g)

		snake.AssertCoordinates(t, gotSnakeCoord, wantSnakeCoord)
		snake.AssertCoordinate(t, *gotFoodCoord, *wantFoodCoord)
	})

	t.Run("should quit game releasing resources", func(t *testing.T) {
		s := snake.NewSnake(width, height)
		cloak := NewStubCloak()
		defer cloak.Stop()

		g := snake.NewGame(s, cloak, fs)
		g.Start(time.Microsecond)

		snake.WaitAndReceiveGameChannels(t, g)
		snake.WaitAndReceiveGameChannels(t, g)

		g.Quit()

		select {
		case cloak.C <- time.Now():
			t.Error("should not be able to send ticks to game")
		case <-time.After(time.Millisecond):
		}
	})
}

func assertNoGameResult(t testing.TB, result *bool) {
	t.Helper()
	if result != nil {
		t.Fatalf("got result %v, want nil", result)
	}
}

func assertNoSnakeCoordinates(t testing.TB, c []snake.Coordinate) {
	t.Helper()
	if c != nil {
		t.Error("shouldn't have got snake coordinates")
	}
}

func assertNoFoodCoordinate(t testing.TB, f *snake.Coordinate) {
	t.Helper()
	if f != nil {
		t.Error("shouldn't have got food coordinate")
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
