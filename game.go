package snake

import "time"

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
}

// NewGame returns a pointer to Game, which handles snake
// methods on cloak ticks
func NewGame(snake *Snake, cloak Cloak, foodProducer FoodGenerator) *Game {
	coordinatesChannel := make(chan []Coordinate)
	movesChannel := make(chan Direction)
	resultChannel := make(chan bool)
	foodChannel := make(chan Coordinate)
	return &Game{snake, cloak, foodProducer, coordinatesChannel, movesChannel, resultChannel, Coordinate{}, foodChannel}
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
	defer close(g.snakeCoordinatesC)
	defer close(g.movesC)
	defer close(g.resultC)
	defer close(g.foodC)
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
