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
	screen tcell.Screen
}

// NewView returns a View struct pointer setting the screen.
func NewView(screen tcell.Screen) *View {
	return &View{screen}
}

// Refresh clears the screen, then prints the snake body on
// the snake coordinates.
func (v *View) Refresh(snakeCoordinates []Coordinate, foodCoordinate Coordinate) {
	v.screen.Clear()
	snakeStyle := tcell.StyleDefault.Foreground(BodyForegroundColor).Background(BodyBackgroundColor)
	for _, c := range snakeCoordinates {
		v.screen.SetContent(c.X, c.Y, BodyRune, nil, snakeStyle)
	}
	foodStyle := tcell.StyleDefault.Foreground(FoodForegroundColor).Background(FoodBackgroundColor)
	v.screen.SetContent(foodCoordinate.X, foodCoordinate.Y, FoodRune, nil, foodStyle)
	v.screen.Show()
}

// Release releases the underlying screen resources.
func (v *View) Release() {
	v.screen.Fini()
}
