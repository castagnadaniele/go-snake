package snake

import "time"

// Game coordinates the snake behaviour with the cloak ticks.
type Game struct {
	snake        Snake
	cloak        Cloak
	coordinatesC chan []Coordinate
	movesC       chan Direction
}

// NewGame returns a pointer to Game, which handles snake
// methods on cloak ticks
func NewGame(snake Snake, cloak Cloak) *Game {
	coordinatesChannel := make(chan []Coordinate)
	movesChannel := make(chan Direction)
	return &Game{snake, cloak, coordinatesChannel, movesChannel}
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
			g.snake.Move(direction)
			g.coordinatesC <- g.snake.GetCoordinates()
		case d := <-g.movesC:
			if g.snake.IsValidMove(d) {
				direction = d
			}
		}
	}
}

// Coordinates returns a []Coordinate receiver channel
// which exposes snake moves after every new tick
func (g *Game) Coordinates() <-chan []Coordinate {
	return g.coordinatesC
}

// SendMove sends d Direction to the internal Direction channel
// which will be pooled inside the Start go routine to change snake direction
func (g *Game) SendMove(d Direction) {
	g.movesC <- d
}
