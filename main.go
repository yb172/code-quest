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
	ClearScreen()
	fmt.Println("\n\033[33;1mGame Starts </>\033[0m\n")
	illusion("A_______________")
}

func game(firstPlace string) bool {
	for i := 0; i < len(firstPlace); i++ {
		if firstPlace[i] == 'A' {
			if i == len(firstPlace)-1 {
				return true
			}
		}
	}
	return false
}

func illusion(firstPlace string) {
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	for {
		fmt.Printf("\r%s", firstPlace)

		if game(firstPlace) {
			break
		}

		firstPlace = keypress(firstPlace)

		time.Sleep(100 * time.Millisecond)
	}

	fmt.Println("\n\033[32mYou won!\033[0m")
}

func keypress(firstPlace string) string {
	_, key, err := keyboard.GetKey()
	if err != nil {
		panic(err)
	}

	if key == keyboard.KeyArrowLeft {
		firstPlace = firstPlace[1:] + "_"
	}

	if key == keyboard.KeyArrowRight {
		firstPlace = "_" + firstPlace[:len(firstPlace)-1]
	}

	if key == keyboard.KeyEsc {
		fmt.Println("\n\033[31mGame Over\033[0m")
		return firstPlace
	}

	return firstPlace
}

func ClearScreen() {
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
