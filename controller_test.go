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

		if !game.Started() {
			t.Errorf("game should have started")
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
}

type GameSpy struct {
	startC            chan struct{}
	SnakeCoordinatesC chan []snake.Coordinate
	FoodCoordinateC   chan snake.Coordinate
	ResultC           chan bool
	MoveC             chan snake.Direction
}

func NewGameSpy() *GameSpy {
	startChannel := make(chan struct{})
	snakeCoordiantesChannel := make(chan []snake.Coordinate)
	foodCoordinatesChannel := make(chan snake.Coordinate)
	resultChannel := make(chan bool)
	moveChannel := make(chan snake.Direction)
	return &GameSpy{
		startC:            startChannel,
		SnakeCoordinatesC: snakeCoordiantesChannel,
		FoodCoordinateC:   foodCoordinatesChannel,
		ResultC:           resultChannel,
		MoveC:             moveChannel,
	}
}

func (g *GameSpy) Started() bool {
	_, ok := <-g.startC
	return !ok
}

func (g *GameSpy) Start(d time.Duration) {
	close(g.startC)
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

type ViewSpy struct {
	DirectionC        chan snake.Direction
	SnakeCoordinatesC chan *[]snake.Coordinate
	FoodCoordinateC   chan *snake.Coordinate
}

func NewViewSpy() *ViewSpy {
	directionChannel := make(chan snake.Direction)
	snakeChannel := make(chan *[]snake.Coordinate)
	foodChannel := make(chan *snake.Coordinate)
	return &ViewSpy{
		DirectionC:        directionChannel,
		SnakeCoordinatesC: snakeChannel,
		FoodCoordinateC:   foodChannel,
	}
}

func (v *ViewSpy) Refresh(snakeCoordinates *[]snake.Coordinate, foodCoordinate *snake.Coordinate) {
	v.SnakeCoordinatesC <- snakeCoordinates
	v.FoodCoordinateC <- foodCoordinate
}

func (v *ViewSpy) ReceiveDirection() <-chan snake.Direction {
	return v.DirectionC
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
