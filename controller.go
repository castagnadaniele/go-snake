package snake

import "time"

// Controller struct coordinates a snake game with a view.
type Controller struct {
	game                GameDirector
	view                ViewHandler
	lastSnakeCoordinate *[]Coordinate
	lastFoodCoordinate  *Coordinate
	gameInterval        time.Duration
	quitC               chan struct{}
}

// NewController returns a Controller pointer initializing the game and the view.
func NewController(game GameDirector, view ViewHandler) *Controller {
	quitChannel := make(chan struct{})
	return &Controller{game, view, nil, nil, 0, quitChannel}
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
	c.gameInterval = d
	c.game.Start(d)
	for {
		select {
		case dir := <-c.view.ReceiveDirection():
			c.game.SendMove(dir)
		case <-c.view.ReceiveNewGameSignal():
			c.game.Restart(c.gameInterval)
		case <-c.view.ReceiveQuitSignal():
			c.game.Quit()
			c.quitC <- struct{}{}
			return
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

// WaitForQuitSignal returns an empty struct receiver channel on which
// the controller sends when it has received a quit signal from view.
// After calling Controller.Start on the main go routine the consumer
// should wait on this channel.
func (c *Controller) WaitForQuitSignal() <-chan struct{} {
	return c.quitC
}
