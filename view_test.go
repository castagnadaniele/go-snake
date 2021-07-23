package snake_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/castagnadaniele/go-snake"
	"github.com/gdamore/tcell/v2"
)

func TestView(t *testing.T) {
	width, height := 60, 60
	snakeCoordinates := &[]snake.Coordinate{{0, 0}, {1, 0}, {2, 0}}

	t.Run("should display snake and food", func(t *testing.T) {
		view, screen := initView(t, width, height)
		defer view.Release()
		foodCoordinate := &snake.Coordinate{6, 6}

		view.Refresh(snakeCoordinates, foodCoordinate)

		for i := 0; i < 3; i++ {
			r, _, s, _ := screen.GetContent(i, 0)
			assertCellRune(t, i, 0, r, snake.BodyRune)
			fg, bg, _ := s.Decompose()
			assertForegroundColor(t, i, 0, fg, snake.BodyForegroundColor)
			assertBackgroundColor(t, i, 0, bg, snake.BodyBackgroundColor)
		}
		r, _, s, _ := screen.GetContent(foodCoordinate.X, foodCoordinate.Y)
		assertCellRune(t, foodCoordinate.X, foodCoordinate.Y, r, snake.FoodRune)
		fg, bg, _ := s.Decompose()
		assertForegroundColor(t, foodCoordinate.X, foodCoordinate.Y, fg, snake.FoodForegroundColor)
		assertBackgroundColor(t, foodCoordinate.X, foodCoordinate.Y, bg, snake.FoodBackgroundColor)
	})

	t.Run("should display snake over food", func(t *testing.T) {
		view, screen := initView(t, width, height)
		defer view.Release()
		foodCoordinate := &snake.Coordinate{0, 0}

		view.Refresh(snakeCoordinates, foodCoordinate)

		r, _, s, _ := screen.GetContent(0, 0)
		assertCellRune(t, 0, 0, r, snake.BodyRune)
		fg, bg, _ := s.Decompose()
		assertForegroundColor(t, 0, 0, fg, snake.BodyForegroundColor)
		assertBackgroundColor(t, 0, 0, bg, snake.BodyBackgroundColor)
	})

	directionTestCases := []struct {
		key tcell.Key
		dir snake.Direction
	}{
		{tcell.KeyUp, snake.Up},
		{tcell.KeyDown, snake.Down},
		{tcell.KeyRight, snake.Right},
		{tcell.KeyLeft, snake.Left},
	}

	for _, tc := range directionTestCases {
		t.Run(fmt.Sprintf("should send %v key and get %v input", tc.key, tc.dir), func(t *testing.T) {
			view, screen := initView(t, width, height)
			defer view.Release()

			screen.InjectKey(tc.key, rune(tc.key), tcell.ModNone)
			direction := <-view.ReceiveDirection()

			snake.AssertDirection(t, direction, tc.dir)
		})
	}

	t.Run("should not display snake if it is nil", func(t *testing.T) {
		view, screen := initView(t, width, height)
		defer view.Release()

		view.Refresh(nil, &snake.Coordinate{0, 0})
		cells, _, _ := screen.GetContents()
		for _, c := range cells {
			if c.Runes[0] == snake.BodyRune {
				t.Fatalf("got snake body rune")
			}
		}
	})

	t.Run("should not display food if it is nil", func(t *testing.T) {
		view, screen := initView(t, width, height)
		defer view.Release()

		view.Refresh(snakeCoordinates, nil)
		cells, _, _ := screen.GetContents()
		for _, c := range cells {
			if c.Runes[0] == snake.FoodRune {
				t.Fatalf("got food rune")
			}
		}
	})

	t.Run("should display win", func(t *testing.T) {
		view, screen := initView(t, width, height)
		defer view.Release()

		view.DisplayWin()
		i := 0
		cells, _, _ := screen.GetContents()
		for _, c := range snake.WinMessage {
			got := cells[i].Runes[0]
			if got != c {
				t.Fatalf("got %c rune, want %c rune", got, c)
			}
			i++
		}
	})

	t.Run("should display lose", func(t *testing.T) {
		view, screen := initView(t, width, height)
		defer view.Release()

		view.DisplayLose()
		i := 0
		cells, _, _ := screen.GetContents()
		for _, c := range snake.LoseMessage {
			got := cells[i].Runes[0]
			if got != c {
				t.Fatalf("got %c rune, want %c rune", got, c)
			}
			i++
		}
	})

	t.Run("should send new game signal on spacebar press", func(t *testing.T) {
		view, screen := initView(t, width, height)
		defer view.Release()

		screen.InjectKey(tcell.KeyRune, ' ', tcell.ModNone)
		select {
		case <-view.ReceiveNewGameSignal():
		case <-time.After(time.Millisecond * 5):
			t.Error("should have received a new game signal")
		}
	})

	t.Run("should send quit game signal on Q press", func(t *testing.T) {
		view, screen := initView(t, width, height)
		defer view.Release()

		keys := []rune{'q', 'Q'}

		for _, r := range keys {
			screen.InjectKey(tcell.KeyRune, r, tcell.ModNone)
			select {
			case <-view.ReceiveQuitSignal():
			case <-time.After(time.Millisecond * 5):
				t.Error("should have received a quit game signal")
			}
		}
	})
}

func assertCellRune(t testing.TB, x, y int, got rune, want rune) {
	t.Helper()
	if got != want {
		t.Errorf("cell at [%d, %d] should contain a %c rune, got %c.", x, y, want, got)
	}
}

func assertForegroundColor(t testing.TB, x, y int, got tcell.Color, want tcell.Color) {
	t.Helper()
	if got != want {
		t.Errorf("got foreground color %v at [%d, %d], want %v", got, x, y, want)
	}
}

func assertBackgroundColor(t testing.TB, x, y int, got tcell.Color, want tcell.Color) {
	t.Helper()
	if got != want {
		t.Errorf("got background color %v at [%d, %d], want %v", got, x, y, want)
	}
}

func initView(t testing.TB, width, height int) (*snake.View, tcell.SimulationScreen) {
	t.Helper()
	screen := tcell.NewSimulationScreen("UTF-8")
	err := screen.Init()
	snake.AssertNoError(t, err)
	screen.SetSize(width, height)
	view := snake.NewView(screen)
	return view, screen
}
