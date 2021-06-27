package snake

type Direction int8

const (
	Up Direction = 1 << iota
	Down
	Left
	Right
)

// Has returns true if flag has d bit turned on
func Has(d, flag Direction) bool { return d&flag != 0 }

func (d Direction) String() string {
	switch d {
	case Up:
		return "Up"
	case Down:
		return "Down"
	case Left:
		return "Left"
	case Right:
		return "Right"
	}
	return "Invalid direction"
}
