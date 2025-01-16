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
	layout    string
	holeIndex int
}

func main() {
	clearScreen()
	worldBuild()
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
		fmt.Printf("\n  %s\n", world.layout)

		if winCheck(world) {
			fmt.Println("\n\033[32mYou won!\033[0m")
			return
		}

		world.layout = keyPress(world)
		time.Sleep(100 * time.Millisecond)
	}
}

func worldBuild() {
	worldSize := flag.Int("world-size", 16, "Specify the size of the world")
	flag.Parse()

	if *worldSize <= 1 {
		fmt.Println("ERR: World size must be greater than 1.")
		os.Exit(1)
	}

	rand.Seed(time.Now().UnixNano())

	holeIndex := rand.Intn(*worldSize-2) + 1

	worldRunes := make([]rune, *worldSize)
	for i := range worldRunes {
		if i == 0 {
			worldRunes[i] = 'A'
		} else if i == holeIndex {
			worldRunes[i] = ' '
		} else {
			worldRunes[i] = '_'
		}
	}

	world := GameWorld{
		layout:    string(worldRunes),
		holeIndex: holeIndex,
	}

	runGame(world)
}

func winCheck(world GameWorld) bool {
	playerIndex := strings.Index(world.layout, "A")
	return playerIndex == world.holeIndex
}

func keyPress(world GameWorld) string {
	_, key, err := keyboard.GetKey()
	if err != nil {
		panic(err)
	}

	index := strings.Index(world.layout, "A")
	if index == -1 {
		panic("Player not found")
	}

	worldRunes := []rune(world.layout)

	switch key {
	case keyboard.KeyArrowLeft:
		if index > 0 && worldRunes[index-1] != ' ' {
			worldRunes[index] = '_'
			worldRunes[index-1] = 'A'
		}
	case keyboard.KeyArrowRight:
		if index < len(worldRunes)-1 && worldRunes[index+1] != ' ' {
			worldRunes[index] = '_'
			worldRunes[index+1] = 'A'
		}
	case keyboard.KeyArrowDown:
		if index == world.holeIndex-1 || index == world.holeIndex+1 {
			worldRunes[index] = '_'
			worldRunes[world.holeIndex] = 'A'
		}
	case keyboard.KeyEsc:
		fmt.Println("\n\033[31mGame Over\033[0m")
		os.Exit(0)
	}

	return string(worldRunes)
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
