package utils

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/yb172/code-quest/models"
)

func WorldBuild() {
	depth := flag.Int("depth", 4, "Number of levels in the game")
	worldSize := flag.Int("world-size", 16, "Size of each level")
	flag.Parse()

	if *depth <= 1 {
		fmt.Println("ERR: Game depth must be greater than 1.")
		os.Exit(1)
	}
	if *worldSize <= 1 {
		fmt.Println("ERR: World size must be greater than 1.")
		os.Exit(1)
	}

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Create and run the game world
	world := createGameWorld(*depth, *worldSize)
	RunGame(world)
}

func createGameWorld(depth, worldSize int) models.GameWorld {
	levels := make([]string, depth)
	holePositions := make([]int, depth)

	// Generate each level with a random hole position
	for level := 0; level < depth; level++ {
		levelRunes := make([]rune, worldSize)
		for i := range levelRunes {
			levelRunes[i] = '_'
		}

		// Mark the player's starting position
		if level == 0 {
			levelRunes[0] = 'A'
		}

		// Randomly place a hole in the level
		holeIndex := rand.Intn(worldSize-2) + 1
		holePositions[level] = holeIndex
		levelRunes[holeIndex] = ' '

		levels[level] = string(levelRunes)
	}

	return models.GameWorld{
		Levels:        levels,
		PlayerPos:     [2]int{0, 0},
		HolePositions: holePositions,
		Depth:         depth,
	}
}

func RunGame(world models.GameWorld) {
	// Open the keyboard for game controls
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		// Ensure keyboard is closed when the game ends
		_ = keyboard.Close()
	}()

	// Game Loop
	for {
		ClearScreen()

		fmt.Println("\n\033[33;1mGame Starts\033[0m")
		fmt.Println("Use ← → to move, ↓ to jump in hole, ESC to exit")

		// Render each level of the game world
		for i, level := range world.Levels {
			if i == world.PlayerPos[0] {
				fmt.Printf("%s\n", level)
			} else {
				fmt.Printf("%s\n", level)
			}
		}

		if winCheck(world) {
			fmt.Println("\n\033[32mCongratulations! You won the game!\033[0m")
			return
		}

		// Wait for keypress and update world
		world = KeyPress(world)

		time.Sleep(100 * time.Millisecond)
	}
}

func winCheck(world models.GameWorld) bool {
	currentLevel := world.PlayerPos[0]
	playerIndex := strings.Index(world.Levels[currentLevel], "A")

	// Player wins if they reach the last level and jump in the hole
	return currentLevel == len(world.Levels)-1 &&
		(playerIndex == world.HolePositions[currentLevel]-1 ||
			playerIndex == world.HolePositions[currentLevel]+1)
}
