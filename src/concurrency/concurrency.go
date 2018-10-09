package main

import (
	"fmt"
	"time"
)

func main() {

	runBoring()

}

func runBoring() {

	// go statement runs the function as usual, but doesn't make the caller wait.
	go boring("boring go routine!!")
	// Running the below function will wait till the loop is completed
	//boring("boring main!")
	time.Sleep(2 * time.Second)
	fmt.Println("You're boring; I'm leaving.")

}

func boring(msg string) {
	for i := 0; i < 5; i++ {
		fmt.Println(msg, i)
		time.Sleep(time.Second)
	}
}
