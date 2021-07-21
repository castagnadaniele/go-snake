package snake_test

import (
	"testing"
	"time"

	"github.com/castagnadaniele/go-snake"
)

func TestController(t *testing.T) {

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
}

type GameSpy struct {
	startC            chan struct{}
	SnakeCoordinatesC chan []snake.Coordinate
	FoodCoordianteC   chan snake.Coordinate
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
		FoodCoordianteC:   foodCoordinatesChannel,
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
	return g.FoodCoordianteC
}

func (g *GameSpy) ReceiveGameResult() <-chan bool {
	return g.ResultC
}

type ViewSpy struct {
	DirectionC           chan snake.Direction
	LastSnakeCoordinates []snake.Coordinate
	LastFoodCoordiante   snake.Coordinate
}

func NewViewSpy() *ViewSpy {
	directionChannel := make(chan snake.Direction)
	return &ViewSpy{DirectionC: directionChannel}
}

func (v *ViewSpy) Refresh(snakeCoordinates []snake.Coordinate, foodCoordinate snake.Coordinate) {
	v.LastSnakeCoordinates = snakeCoordinates
	v.LastFoodCoordiante = foodCoordinate
}

func (v *ViewSpy) ReceiveDirection() <-chan snake.Direction {
	return v.DirectionC
}
