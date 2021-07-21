package snake

import "time"

// Controller struct coordinates a snake game with a view.
type Controller struct {
	game GameDirector
	view ViewHandler
}

// NewController returns a Controller pointer initializing the game and the view.
func NewController(game GameDirector, view ViewHandler) *Controller {
	return &Controller{game, view}
}

// Start starts the controller internal game, then loops and waits on the view
// direction channel.
//
// Should be used as a go routine.
func (c *Controller) Start(d time.Duration) {
	c.game.Start(d)
	for dir := range c.view.ReceiveDirection() {
		c.game.SendMove(dir)
	}
}
