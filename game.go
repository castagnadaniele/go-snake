package snake

import "time"

// Game coordinates the snake behaviour with the cloak ticks.
type Game struct {
	snake          *Snake
	cloak          Cloak
	foodProducer   FoodGenerator
	coordinatesC   chan []Coordinate
	movesC         chan Direction
	resultC        chan bool
	foodCoordinate Coordinate
}

// NewGame returns a pointer to Game, which handles snake
// methods on cloak ticks
func NewGame(snake *Snake, cloak Cloak, foodProducer FoodGenerator) (*Game, error) {
	coordinatesChannel := make(chan []Coordinate)
	movesChannel := make(chan Direction)
	resultChannel := make(chan bool)
	food, err := foodProducer.Generate(snake.GetCoordinates())
	if err != nil {
		return nil, err
	}
	return &Game{snake, cloak, foodProducer, coordinatesChannel, movesChannel, resultChannel, food}, nil
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
	defer close(g.coordinatesC)
	defer close(g.movesC)
	defer close(g.resultC)
	direction := g.snake.Face()
	for {
		select {
		case <-g.cloak.Tick():
			result := g.handleMove(direction)
			if !result {
				g.resultC <- false
				return
			}
		case d := <-g.movesC:
			if g.snake.IsValidMove(d) {
				direction = d
			}
		}
	}
}

func (g *Game) handleMove(d Direction) bool {
	err := g.snake.Move(d)
	if err == ErrHeadOutOfBoard {
		return false
	}
	coord := g.snake.GetCoordinates()
	head := coord[0]
	if head.X == g.foodCoordinate.X && head.Y == g.foodCoordinate.Y {
		err = g.snake.Grow()
		coord = g.snake.GetCoordinates()
		if err != nil {
			return false
		}
		g.foodCoordinate, err = g.foodProducer.Generate(coord)
		if err != nil {
			return false
		}
	}
	g.coordinatesC <- coord
	return true
}

// SendMove sends d Direction to the internal Direction channel
// which will be pooled inside the Start go routine to change snake direction
func (g *Game) SendMove(d Direction) {
	g.movesC <- d
}

// ReceiveResult returns a ([]Coordinate, *bool) tuple with snake moves
// or game result. It waits to receive values from the game internal channels.
func (g *Game) ReceiveResult() (coordinate []Coordinate, result *bool) {
	select {
	case c := <-g.coordinatesC:
		return c, nil
	case r := <-g.resultC:
		return nil, &r
	}
}
