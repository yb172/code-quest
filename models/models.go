package models

type GameWorld struct {
	Levels        []string
	PlayerPos     [2]int
	HolePositions []int
	Depth         int
}
