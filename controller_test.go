package snake_test

import (
	"testing"
	"time"

	"github.com/castagnadaniele/go-snake"
)

func TestController(t *testing.T) {
	foodCoordinate := snake.Coordinate{0, 0}
	snakeCoordinates := []snake.Coordinate{{0, 1}}

	t.Run("should start game", func(t *testing.T) {
		view := NewViewSpy()
		game := NewGameSpy()
		controller := snake.NewController(game, view)

		go controller.Start(time.Microsecond)

		select {
		case <-game.StartC:
		case <-time.After(time.Millisecond * 5):
			t.Error("game should have started")
		}
	})

	t.Run("should send move from view to game", func(t *testing.T) {
		view := NewViewSpy()
		game := NewGameSpy()
		controller := snake.NewController(game, view)

		go controller.Start(time.Microsecond)

		view.DirectionC <- snake.Up

		got := <-game.MoveC
		snake.AssertDirection(t, got, snake.Up)
	})

	t.Run("should refresh view when game sends snake coordinates", func(t *testing.T) {
		view := NewViewSpy()
		game := NewGameSpy()
		controller := snake.NewController(game, view)

		go controller.Start(time.Microsecond)

		game.SnakeCoordinatesC <- snakeCoordinates
		gotSnake := <-view.SnakeCoordinatesC
		gotFood := <-view.FoodCoordinateC
		assertCoordinatesNotNil(t, gotSnake)
		assertCoordinateNil(t, gotFood)
		snake.AssertCoordinates(t, *gotSnake, snakeCoordinates)
	})

	t.Run("should refresh view when game send food coordinate", func(t *testing.T) {
		view := NewViewSpy()
		game := NewGameSpy()
		controller := snake.NewController(game, view)

		go controller.Start(time.Microsecond)

		game.FoodCoordinateC <- foodCoordinate
		gotSnake := <-view.SnakeCoordinatesC
		gotFood := <-view.FoodCoordinateC
		assertCoordinatesNil(t, gotSnake)
		assertCoordinateNotNil(t, gotFood)
		snake.AssertCoordinate(t, *gotFood, foodCoordinate)
	})

	t.Run("should refresh view when game sends snake coordinate with last food sent", func(t *testing.T) {
		view := NewViewSpy()
		game := NewGameSpy()
		controller := snake.NewController(game, view)

		go controller.Start(time.Microsecond)

		game.FoodCoordinateC <- foodCoordinate
		<-view.SnakeCoordinatesC
		<-view.FoodCoordinateC
		game.SnakeCoordinatesC <- snakeCoordinates
		gotSnake := <-view.SnakeCoordinatesC
		gotFood := <-view.FoodCoordinateC
		assertCoordinatesNotNil(t, gotSnake)
		assertCoordinateNotNil(t, gotFood)
		snake.AssertCoordinates(t, *gotSnake, snakeCoordinates)
		snake.AssertCoordinate(t, *gotFood, foodCoordinate)
	})

	t.Run("should refresh view when game sends food coordinate with last snake sent", func(t *testing.T) {
		view := NewViewSpy()
		game := NewGameSpy()
		controller := snake.NewController(game, view)

		go controller.Start(time.Microsecond)

		game.SnakeCoordinatesC <- snakeCoordinates
		<-view.SnakeCoordinatesC
		<-view.FoodCoordinateC
		game.FoodCoordinateC <- foodCoordinate
		gotSnake := <-view.SnakeCoordinatesC
		gotFood := <-view.FoodCoordinateC
		assertCoordinatesNotNil(t, gotSnake)
		assertCoordinateNotNil(t, gotFood)
		snake.AssertCoordinates(t, *gotSnake, snakeCoordinates)
		snake.AssertCoordinate(t, *gotFood, foodCoordinate)
	})

	t.Run("should display win when game send win result", func(t *testing.T) {
		view := NewViewSpy()
		game := NewGameSpy()
		controller := snake.NewController(game, view)

		go controller.Start(time.Microsecond)

		game.ResultC <- true
		select {
		case <-view.WinC:
		case <-view.LoseC:
			t.Error("view should have displayed a win instead of a lose")
		case <-time.After(time.Millisecond * 5):
			t.Error("view should have displayed a win")
		}
	})

	t.Run("should display lose when game send lose result", func(t *testing.T) {
		view := NewViewSpy()
		game := NewGameSpy()
		controller := snake.NewController(game, view)

		go controller.Start(time.Microsecond)

		game.ResultC <- false
		select {
		case <-view.LoseC:
		case <-view.WinC:
			t.Error("view should have diplayed a lose instead of a win")
		case <-time.After(time.Millisecond * 5):
			t.Error("view should have displayed a lose")
		}
	})

	t.Run("should restart game when receives new game signal from view", func(t *testing.T) {
		view := NewViewSpy()
		game := NewGameSpy()

		controller := snake.NewController(game, view)

		want := time.Microsecond
		go controller.Start(want)

		select {
		case view.NewGameC <- struct{}{}:
		case <-time.After(time.Millisecond * 5):
			t.Errorf("view should have received a new game signal")
		}

		select {
		case got := <-game.RestartC:
			if got != want {
				t.Errorf("got duration %v, want duration %v", got, want)
			}
		case <-time.After(time.Millisecond * 5):
			t.Errorf("game should have restarted")
		}
	})
}

type GameSpy struct {
	StartC            chan struct{}
	SnakeCoordinatesC chan []snake.Coordinate
	FoodCoordinateC   chan snake.Coordinate
	ResultC           chan bool
	MoveC             chan snake.Direction
	RestartC          chan time.Duration
}

func NewGameSpy() *GameSpy {
	startChannel := make(chan struct{}, 1)
	snakeCoordiantesChannel := make(chan []snake.Coordinate)
	foodCoordinatesChannel := make(chan snake.Coordinate)
	resultChannel := make(chan bool)
	moveChannel := make(chan snake.Direction)
	restartChannel := make(chan time.Duration)
	return &GameSpy{
		StartC:            startChannel,
		SnakeCoordinatesC: snakeCoordiantesChannel,
		FoodCoordinateC:   foodCoordinatesChannel,
		ResultC:           resultChannel,
		MoveC:             moveChannel,
		RestartC:          restartChannel,
	}
}

func (g *GameSpy) Start(d time.Duration) {
	g.StartC <- struct{}{}
}

func (g *GameSpy) SendMove(d snake.Direction) {
	g.MoveC <- d
}

func (g *GameSpy) ReceiveSnakeCoordinates() <-chan []snake.Coordinate {
	return g.SnakeCoordinatesC
}

func (g *GameSpy) ReceiveFoodCoordinate() <-chan snake.Coordinate {
	return g.FoodCoordinateC
}

func (g *GameSpy) ReceiveGameResult() <-chan bool {
	return g.ResultC
}

func (g *GameSpy) Restart(d time.Duration) {
	g.RestartC <- d
}

type ViewSpy struct {
	DirectionC        chan snake.Direction
	SnakeCoordinatesC chan *[]snake.Coordinate
	FoodCoordinateC   chan *snake.Coordinate
	WinC              chan struct{}
	LoseC             chan struct{}
	NewGameC          chan struct{}
}

func NewViewSpy() *ViewSpy {
	directionChannel := make(chan snake.Direction)
	snakeChannel := make(chan *[]snake.Coordinate)
	foodChannel := make(chan *snake.Coordinate)
	winChannel := make(chan struct{})
	loseChannel := make(chan struct{})
	newGameChannel := make(chan struct{})
	return &ViewSpy{
		DirectionC:        directionChannel,
		SnakeCoordinatesC: snakeChannel,
		FoodCoordinateC:   foodChannel,
		WinC:              winChannel,
		LoseC:             loseChannel,
		NewGameC:          newGameChannel,
	}
}

func (v *ViewSpy) Refresh(snakeCoordinates *[]snake.Coordinate, foodCoordinate *snake.Coordinate) {
	v.SnakeCoordinatesC <- snakeCoordinates
	v.FoodCoordinateC <- foodCoordinate
}

func (v *ViewSpy) ReceiveDirection() <-chan snake.Direction {
	return v.DirectionC
}

func (v *ViewSpy) DisplayWin() {
	v.WinC <- struct{}{}
}

func (v *ViewSpy) DisplayLose() {
	v.LoseC <- struct{}{}
}

func (v *ViewSpy) ReceiveNewGameSignal() <-chan struct{} {
	return v.NewGameC
}

func assertCoordinateNil(t testing.TB, got *snake.Coordinate) {
	t.Helper()
	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

func assertCoordinateNotNil(t testing.TB, got *snake.Coordinate) {
	t.Helper()
	if got == nil {
		t.Error("should have not got nil")
	}
}

func assertCoordinatesNil(t testing.TB, got *[]snake.Coordinate) {
	t.Helper()
	if got != nil {
		t.Errorf("got %v, want nil", got)
	}
}

func assertCoordinatesNotNil(t testing.TB, got *[]snake.Coordinate) {
	t.Helper()
	if got == nil {
		t.Error("should have not got nil")
	}
}
