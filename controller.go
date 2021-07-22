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
// direction channel, on the game snake coordinates receiver channel, on
// the game food coordinate receiver channel and on the game result receiver channel.
// When it receives a new direction from the view it sends it to the game.
// When it receives new snake or food coordinates it refreshes the view screen.
// When it receives a game result it display win or lose accordingly to the result.
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
		case r := <-c.game.ReceiveGameResult():
			if r {
				c.view.DisplayWin()
			} else {
				c.view.DisplayLose()
			}
		}
	}
}
