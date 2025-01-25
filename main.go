package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

type GameWorld struct {
	levels        []string
	playerPos     [2]int
	holePositions []int
	depth         int
}

func main() {
	clearScreen()
	worldBuild()
}

func worldBuild() {
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

	rand.Seed(time.Now().UnixNano())

	world := createGameWorld(*depth, *worldSize)
	runGame(world)
}

func createGameWorld(depth, worldSize int) GameWorld {
	levels := make([]string, depth)
	holePositions := make([]int, depth)

	for level := 0; level < depth; level++ {
		levelRunes := make([]rune, worldSize)
		for i := range levelRunes {
			levelRunes[i] = '_'
		}

		if level == 0 {
			levelRunes[0] = 'A'
		}

		holeIndex := rand.Intn(worldSize-2) + 1
		holePositions[level] = holeIndex
		levelRunes[holeIndex] = ' '

		levels[level] = string(levelRunes)
	}

	return GameWorld{
		levels:        levels,
		playerPos:     [2]int{0, 0},
		holePositions: holePositions,
		depth:         depth,
	}
}

func runGame(world GameWorld) {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		clearScreen()

		fmt.Println("\n\033[33;1mGame Starts\033[0m")
		fmt.Println("Use ← → to move, ↓ to jump in hole, ESC to exit")

		for i, level := range world.levels {
			if i == world.playerPos[0] {
				fmt.Printf("%s\n", level)
			} else {
				fmt.Printf("%s\n", level)
			}
		}

		if winCheck(world) {
			fmt.Println("\n\033[32mCongratulations! You won the game!\033[0m")
			return
		}

		world = keyPress(world)

		time.Sleep(100 * time.Millisecond)
	}
}

func winCheck(world GameWorld) bool {
	currentLevel := world.playerPos[0]
	playerIndex := strings.Index(world.levels[currentLevel], "A")

	return currentLevel == len(world.levels)-1 &&
		(playerIndex == world.holePositions[currentLevel]-1 ||
			playerIndex == world.holePositions[currentLevel]+1)
}

func keyPress(world GameWorld) GameWorld {
	_, key, err := keyboard.GetKey()
	if err != nil {
		panic(err)
	}

	currentLevel := world.playerPos[0]
	levels := make([]string, len(world.levels))
	copy(levels, world.levels)

	index := strings.Index(levels[currentLevel], "A")
	if index == -1 {
		panic("Player not found")
	}

	levelRunes := []rune(levels[currentLevel])

	switch key {
	case keyboard.KeyArrowLeft:
		if index > 0 && levelRunes[index-1] != ' ' {
			levelRunes[index] = '_'
			levelRunes[index-1] = 'A'
			world.playerPos[1] = index - 1
		}
	case keyboard.KeyArrowRight:
		if index < len(levelRunes)-1 && levelRunes[index+1] != ' ' {
			levelRunes[index] = '_'
			levelRunes[index+1] = 'A'
			world.playerPos[1] = index + 1
		}
	case keyboard.KeyArrowDown:
		holePos := world.holePositions[currentLevel]
		if currentLevel < len(world.levels)-1 {
			if index == holePos-1 || index == holePos+1 {
				levelRunes[index] = '_'
				levels[currentLevel] = string(levelRunes)

				nextLevel := currentLevel + 1
				nextLevelRunes := []rune(levels[nextLevel])

				nextLevelRunes[holePos] = 'A'
				world.playerPos = [2]int{nextLevel, holePos}

				levels[nextLevel] = string(nextLevelRunes)
				world.levels = levels
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

	levels[currentLevel] = string(levelRunes)
	world.levels = levels

	return world
}

func clearScreen() {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "cls")
	} else {
		cmd = exec.Command("clear")
	}
	cmd.Stdout = os.Stdout
	_ = cmd.Run()
}
