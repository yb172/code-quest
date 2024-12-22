package main

import (
	"fmt"
	"os"
)

func main() {
	var name string
	fmt.Println("Hi, what is your name?")
	fmt.Fscan(os.Stdin, &name)
	fmt.Printf("Hi, %s!\n", name)
}
