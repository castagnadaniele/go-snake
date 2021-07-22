package snake

import "github.com/gdamore/tcell/v2"

const BodyRune = '▮'
const BodyForegroundColor = tcell.ColorWhite
const BodyBackgroundColor = tcell.ColorGray
const FoodRune = '◆'
const FoodForegroundColor = tcell.ColorRed
const FoodBackgroundColor = tcell.ColorBlack
const WinMessage = "Game won! Press SPACEBAR to start a new game or press Q to quit..."
const LoseMessage = "Game lost! Press SPACEBAR to start a new game or press Q to quit..."

// ViewHandler interface defines how a view should handle
// screen refresh and how should expose snake's change direction input.
type ViewHandler interface {
	// Refresh should receive the snake and food coordinates and should display them.
	Refresh(snakeCoordinates *[]Coordinate, foodCoordinate *Coordinate)
	// ReceiveDirection should return a Direction receiver channel on which the ViewHandler
	// should send new change direction input from the user.
	ReceiveDirection() <-chan Direction
	// DisplayWin should display a win screen.
	DisplayWin()
	// DisplayLose should display a lose screen.
	DisplayLose()
}

// View struct which prints the snake game elements on terminal.
type View struct {
	screen      tcell.Screen
	directionC  chan Direction
	eventsC     chan tcell.Event
	quitEventsC chan struct{}
}

// NewView returns a View struct pointer setting the screen,
// starting the screen events loop channel in a go routine
// and starts polling the screen events channel for directions
// in another go routine.
func NewView(screen tcell.Screen) *View {
	directionChannel := make(chan Direction)
	eventsChannel := make(chan tcell.Event)
	quitEventsC := make(chan struct{})
	go screen.ChannelEvents(eventsChannel, quitEventsC)
	view := &View{screen, directionChannel, eventsChannel, quitEventsC}
	go view.pollDirectionKeys()
	return view
}

// Refresh clears the screen, then prints the snake body on
// the snake coordinates and the food on the food coordinates.
// The snake body will be printed overwriting the food, if their coordinates overlap.
// It will not print the respective coordinates if the snake or the food coordinates are nil.
func (v *View) Refresh(snakeCoordinates *[]Coordinate, foodCoordinate *Coordinate) {
	if snakeCoordinates == nil && foodCoordinate == nil {
		return
	}
	v.screen.Clear()
	if foodCoordinate != nil {
		foodStyle := tcell.StyleDefault.Foreground(FoodForegroundColor).Background(FoodBackgroundColor)
		v.screen.SetContent(foodCoordinate.X, foodCoordinate.Y, FoodRune, nil, foodStyle)
	}
	if snakeCoordinates != nil {
		snakeStyle := tcell.StyleDefault.Foreground(BodyForegroundColor).Background(BodyBackgroundColor)
		for _, c := range *snakeCoordinates {
			v.screen.SetContent(c.X, c.Y, BodyRune, nil, snakeStyle)
		}
	}
	v.screen.Show()
}

// Release releases the underlying screen resources.
func (v *View) Release() {
	// Screen.ChannelEvents will close v.eventsC after we close v.quitEventsC
	close(v.quitEventsC)
	close(v.directionC)
	v.screen.Fini()
}

// ReceiveDirection returns a Direction receiver channel
// which will be fed when the screen will receive
// directional key events.
func (v *View) ReceiveDirection() <-chan Direction {
	return v.directionC
}

// DisplayWin clears the screen and displays a win message.
func (v *View) DisplayWin() {
	v.printMessage(WinMessage)
}

// DisplayLose clears the screen and displays a lose message.
func (v *View) DisplayLose() {
	v.printMessage(LoseMessage)
}

func (v *View) pollDirectionKeys() {
	for e := range v.eventsC {
		if keyEvent, ok := e.(*tcell.EventKey); ok {
			switch keyEvent.Key() {
			case tcell.KeyUp:
				v.directionC <- Up
			case tcell.KeyDown:
				v.directionC <- Down
			case tcell.KeyRight:
				v.directionC <- Right
			case tcell.KeyLeft:
				v.directionC <- Left
			}
		}
	}
}

func (v *View) printMessage(message string) {
	v.screen.Clear()
	width, _ := v.screen.Size()
	x, y := 0, 0
	for _, c := range message {
		if x >= width {
			x = 0
			y++
		}
		v.screen.SetContent(x, y, c, nil, tcell.StyleDefault)
		x++
	}
	v.screen.Show()
}
