package snake

// Coordinate implements board coordinates.
type Coordinate struct {
	X int
	Y int
}

func contains(arr []Coordinate, c Coordinate) bool {
	for _, item := range arr {
		if item.X == c.X && item.Y == c.Y {
			return true
		}
	}
	return false
}
