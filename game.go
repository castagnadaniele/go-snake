package snake

import "time"

// Game coordinates the snake behaviour with the cloak ticks.
type Game struct {
	snake        Snake
	cloak        Cloak
	coordinatesC chan []Coordinate
}

// NewGame returns a pointer to Game, which handles snake
// methods on cloak ticks
func NewGame(snake Snake, cloak Cloak) *Game {
	coordinatesChannel := make(chan []Coordinate)
	return &Game{snake, cloak, coordinatesChannel}
}

// Start starts cloak to tick every d time.Duration,
// then starts a go routine to loop on the ticker events
// moving the snake and sending the new coordinates on
// the internal channel
func (g *Game) Start(d time.Duration) {
	g.cloak.Start(d)
	go func() {
		defer close(g.coordinatesC)
		for range g.cloak.Tick() {
			g.snake.Move(Left)
			g.coordinatesC <- g.snake.GetCoordinates()
		}
	}()
}

// Coordinates returns a []Coordinate receiver channel
// which exposes snake moves after every new tick
func (g *Game) Coordinates() <-chan []Coordinate {
	return g.coordinatesC
}
