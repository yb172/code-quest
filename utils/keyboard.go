package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/eiannone/keyboard"
	"github.com/yb172/code-quest/models"
)

/*
	Whole this file ( keyboard.go ) was made by looking at the documentation of the package
`github.com/eiannone/keyboard` to be better in understanding the code from this package
 visit the link: https://github.com/eiannone/keyboard`, it has a good documentation there,
 but it is not made in GitBook, which is not good :(
    ~ Arsenii Trutnev
*/

func KeyPress(world models.GameWorld) models.GameWorld {
	// Get the key pressed by the user
	_, key, err := keyboard.GetKey()
	if err != nil {
		panic(err)
	}

	// Get the current level of the player
	currentLevel := world.PlayerPos[0]
	// Create a copy of the levels to avoid modifying the original
	levels := make([]string, len(world.Levels))
	copy(levels, world.Levels)

	// Find the player's current position (`A`) in the current level
	index := strings.Index(levels[currentLevel], "A")
	if index == -1 {
		panic("Player not found")
	}

	// Convert the level string to a slice of runes to modify it
	levelRunes := []rune(levels[currentLevel])

	// Process the key press and update the player's position
	switch key {

	// Move left
	case keyboard.KeyArrowLeft:
		if index > 0 && levelRunes[index-1] != ' ' {
			levelRunes[index] = '_'
			levelRunes[index-1] = 'A'
			world.PlayerPos[1] = index - 1
		}
		// Move right
	case keyboard.KeyArrowRight:
		if index < len(levelRunes)-1 && levelRunes[index+1] != ' ' {
			levelRunes[index] = '_'
			levelRunes[index+1] = 'A'
			world.PlayerPos[1] = index + 1
		}
		// Jumping into a hole
	case keyboard.KeyArrowDown:
		holePos := world.HolePositions[currentLevel]
		if currentLevel < len(world.Levels)-1 {
			if index == holePos-1 || index == holePos+1 {
				levelRunes[index] = '_'
				levels[currentLevel] = string(levelRunes)

				nextLevel := currentLevel + 1
				nextLevelRunes := []rune(levels[nextLevel])

				nextLevelRunes[holePos] = 'A'
				world.PlayerPos = [2]int{nextLevel, holePos}

				levels[nextLevel] = string(nextLevelRunes)
				world.Levels = levels
			}
		} else {
			if index == holePos-1 || index == holePos+1 {
				fmt.Println("\n\033[32mCongratulations! You won the game!\033[0m")
				os.Exit(0)
			}
		}
	case keyboard.KeyEsc:
		fmt.Println("\n\033[31mGame Over\033[0m")
		os.Exit(0)
	}

	// Update the level with the new player position
	levels[currentLevel] = string(levelRunes)
	world.Levels = levels

	// Return the updated world
	return world
}
