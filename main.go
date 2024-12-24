package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
)

func main() {
	clearScreen()
	worldBuild()
}

func runGame(world string) {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		clearScreen()
		fmt.Println("\n\033[33;1mGame Starts\033[0m")
		fmt.Printf("\n  %s\n", world)

		if winCheck(world) {
			fmt.Println("\n\033[32mYou won!\033[0m")
			return
		}

		world = keyPress(world)
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

	world := "A" + strings.Repeat("_", *worldSize-1)
	runGame(world)
}

func winCheck(world string) bool {
	return world[len(world)-1] == 'A'
}

func keyPress(world string) string {
	_, key, err := keyboard.GetKey()
	if err != nil {
		panic(err)
	}

	index := strings.Index(world, "A")
	if index == -1 {
		panic("Player not found")
	}

	worldRunes := []rune(world)

	switch key {
	case keyboard.KeyArrowLeft:
		if index > 0 {
			worldRunes[index] = '_'
			worldRunes[index-1] = 'A'
		}
	case keyboard.KeyArrowRight:
		if index < len(worldRunes)-1 {
			worldRunes[index] = '_'
			worldRunes[index+1] = 'A'
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
