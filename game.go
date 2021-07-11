package snake

import "time"

// Game coordinates the snake behaviour with the cloak ticks.
type Game struct {
	snake        *Snake
	cloak        Cloak
	coordinatesC chan []Coordinate
	movesC       chan Direction
	resultC      chan bool
}

// NewGame returns a pointer to Game, which handles snake
// methods on cloak ticks
func NewGame(snake *Snake, cloak Cloak) *Game {
	coordinatesChannel := make(chan []Coordinate)
	movesChannel := make(chan Direction)
	resultChannel := make(chan bool)
	return &Game{snake, cloak, coordinatesChannel, movesChannel, resultChannel}
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
	direction := g.snake.Face()
	for {
		select {
		case <-g.cloak.Tick():
			err := g.snake.Move(direction)
			if err == ErrHeadOutOfBoard {
				g.resultC <- false
				return
			}
			g.coordinatesC <- g.snake.GetCoordinates()
		case d := <-g.movesC:
			if g.snake.IsValidMove(d) {
				direction = d
			}
		}
	}
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
