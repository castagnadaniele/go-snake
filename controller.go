package snake

import "time"

// Controller struct coordinates a snake game with a view.
type Controller struct {
	game                GameDirector
	view                ViewHandler
	lastSnakeCoordinate *[]Coordinate
	lastFoodCoordinate  *Coordinate
}

// NewController returns a Controller pointer initializing the game and the view.
func NewController(game GameDirector, view ViewHandler) *Controller {
	return &Controller{game, view, nil, nil}
}

// Start starts the controller internal game, then loops and waits on the view
// direction channel, on the game snake coordinates receiver channel and on
// the game food coordinate receiver channel.
//
// Should be used as a go routine.
func (c *Controller) Start(d time.Duration) {
	c.game.Start(d)
	for {
		select {
		case dir := <-c.view.ReceiveDirection():
			c.game.SendMove(dir)
		case sc := <-c.game.ReceiveSnakeCoordinates():
			c.lastSnakeCoordinate = &sc
			c.view.Refresh(c.lastSnakeCoordinate, c.lastFoodCoordinate)
		case fc := <-c.game.ReceiveFoodCoordinate():
			c.lastFoodCoordinate = &fc
			c.view.Refresh(c.lastSnakeCoordinate, c.lastFoodCoordinate)
		}
	}
}
