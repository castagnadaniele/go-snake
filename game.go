package snake

import "time"

// GameDirector interface defines how to coordinate the snake and food
// interaction in a game.
type GameDirector interface {
	// Start should exec a new snake move after an interval, should handle
	// snake collision with food, should handle new food generation and should
	// change snake face direction on user input.
	//
	// Start should run a go routine which does the above operations.
	Start(d time.Duration)
	// SendMove should send the new direction in an internal channel.
	SendMove(d Direction)
	// ReceiveSnakeCoordinates should expose a receiver channel which emits
	// the new snake coordinates after each interval.
	ReceiveSnakeCoordinates() <-chan []Coordinate
	// ReceiveFoodCoordinate should expose a receiver channel which emits
	// the new food coordinate after the snake eats the food.
	ReceiveFoodCoordinate() <-chan Coordinate
	// ReceiveGameResult should expose a receiver channel which emits
	// when the game is won or is lost.
	ReceiveGameResult() <-chan bool
}

// Game coordinates the snake behaviour with the cloak ticks.
type Game struct {
	snake             *Snake
	cloak             Cloak
	foodProducer      FoodGenerator
	snakeCoordinatesC chan []Coordinate
	movesC            chan Direction
	resultC           chan bool
	foodCoordinate    Coordinate
	foodC             chan Coordinate
	quitEventRoutineC chan struct{}
}

// NewGame returns a pointer to Game, which handles snake
// methods on cloak ticks
func NewGame(snake *Snake, cloak Cloak, foodProducer FoodGenerator) *Game {
	coordinatesChannel := make(chan []Coordinate)
	movesChannel := make(chan Direction)
	resultChannel := make(chan bool)
	foodChannel := make(chan Coordinate)
	quitEventRoutineChannel := make(chan struct{})
	return &Game{
		snake,
		cloak,
		foodProducer,
		coordinatesChannel,
		movesChannel,
		resultChannel,
		Coordinate{},
		foodChannel,
		quitEventRoutineChannel,
	}
}

// Start starts cloak to tick every d time.Duration,
// then starts a go routine to loop on the ticker events
// moving the snake and sending the new coordinates on
// the internal channel
func (g *Game) Start(d time.Duration) {
	g.cloak.Start(d)
	go g.eventRoutine()
}

func (g *Game) eventRoutine() {
	// defer close(g.snakeCoordinatesC)
	// defer close(g.movesC)
	// defer close(g.resultC)
	// defer close(g.foodC)
	direction := g.snake.Face()
	g.sendInitSnakeAndFoodCoordinates()
	for {
		select {
		case <-g.cloak.Tick():
			result := g.handleMove(direction)
			if result != nil {
				g.resultC <- *result
				return
			}
		case d := <-g.movesC:
			if g.snake.IsValidMove(d) {
				direction = d
			}
		case <-g.quitEventRoutineC:
			return
		}
	}
}

func (g *Game) sendInitSnakeAndFoodCoordinates() {
	g.snakeCoordinatesC <- g.snake.GetCoordinates()
	var err error
	g.foodCoordinate, err = g.foodProducer.Generate(g.snake.GetCoordinates())
	if err != nil {
		panic(err)
	}
	g.foodC <- g.foodCoordinate
}

func (g *Game) handleMove(d Direction) *bool {
	result := false
	err := g.snake.Move(d)
	if err == ErrHeadOutOfBoard {
		return &result
	}
	coord := g.snake.GetCoordinates()
	head := coord[0]
	if head.X == g.foodCoordinate.X && head.Y == g.foodCoordinate.Y {
		err = g.snake.Grow()
		if err != nil {
			return &result
		}
		coord = g.snake.GetCoordinates()
		g.snakeCoordinatesC <- coord
		g.foodCoordinate, err = g.foodProducer.Generate(coord)
		if err != nil {
			result = true
			return &result
		}
		g.foodC <- g.foodCoordinate
		return nil
	}
	g.snakeCoordinatesC <- coord
	return nil
}

// SendMove sends d Direction to the internal Direction channel
// which will be pooled inside the Start go routine to change snake direction
func (g *Game) SendMove(d Direction) {
	g.movesC <- d
}

// ReceiveSnakeCoordinates returns the snake coordinates receive channel.
func (g *Game) ReceiveSnakeCoordinates() <-chan []Coordinate {
	return g.snakeCoordinatesC
}

// ReceiveFoodCoordinate returns the food coordinate receive channel.
func (g *Game) ReceiveFoodCoordinate() <-chan Coordinate {
	return g.foodC
}

// ReceiveGameResult returns the game result receive channel.
func (g *Game) ReceiveGameResult() <-chan bool {
	return g.resultC
}

// Restart stops the game internal go routine, reset the snake and starts
// a new game event loop internal go routine.
func (g *Game) Restart(d time.Duration) {
	g.quitEventRoutineC <- struct{}{}
	g.snake.Reset()
	go g.eventRoutine()
}
