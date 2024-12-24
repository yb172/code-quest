package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/eiannone/keyboard"
)

func main() {
	clearScreen()
	fmt.Println("\n\033[33;1mGame Starts </>\033[0m\n")
	runGame("A_______________")
}

func runGame(world string) {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		fmt.Printf("\r%s", world)

		if winCheck(world) {
			break
		}

		world = keyPress(world)

		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\n\033[32mYou won!\033[0m")
}

func winCheck(world string) bool {
	for i := 0; i < len(world); i++ {
		if world[i] == 'A' {
			if i == len(world)-1 {
				return true
			}
		}
	}
	return false
}

func keyPress(world string) string {
	_, key, err := keyboard.GetKey()
	if err != nil {
		panic(err)
	}

	if key == keyboard.KeyArrowLeft {
		world = world[1:] + "_"
	}

	if key == keyboard.KeyArrowRight {
		world = "_" + world[:len(world)-1]
	}

	if key == keyboard.KeyEsc {
		fmt.Println("\n\033[31mGame Over\033[0m")
		return world
	}

	return world
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
