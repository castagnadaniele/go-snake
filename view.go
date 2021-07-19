package snake

import "github.com/gdamore/tcell/v2"

const BodyRune = '▮'
const BodyForegroundColor = tcell.ColorWhite
const BodyBackgroundColor = tcell.ColorGray
const FoodRune = '◆'
const FoodForegroundColor = tcell.ColorRed
const FoodBackgroundColor = tcell.ColorBlack

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
func (v *View) Refresh(snakeCoordinates []Coordinate, foodCoordinate Coordinate) {
	v.screen.Clear()
	foodStyle := tcell.StyleDefault.Foreground(FoodForegroundColor).Background(FoodBackgroundColor)
	v.screen.SetContent(foodCoordinate.X, foodCoordinate.Y, FoodRune, nil, foodStyle)
	snakeStyle := tcell.StyleDefault.Foreground(BodyForegroundColor).Background(BodyBackgroundColor)
	for _, c := range snakeCoordinates {
		v.screen.SetContent(c.X, c.Y, BodyRune, nil, snakeStyle)
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