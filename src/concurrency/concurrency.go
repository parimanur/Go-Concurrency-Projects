package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	fmt.Println("Started runBoring")
	runBoring()
	fmt.Println("Finished runBoring\n\n")

	time.Sleep(5 * time.Second)

	fmt.Println("Starting runBoringWithChan")
	runBoringWithChan()
	fmt.Println("Finished runBoringWithChan\n\n")
}

func runBoringWithChan() {
	c := make(chan string)
	go boringWithChan("boring with chan c!", c)
	for i := 0; i < 5; i++ {
		fmt.Printf("You say: %q\n", <-c) // wait for a value to be sent.
	}
	fmt.Println("You're boring; I'm leaving.")
}

func runBoring() {

	// go statement runs the function as usual, but doesn't make the caller wait.
	go boring("boring go routine!!")
	// Running the below function will wait till the loop is completed
	//boring("boring main!")
	time.Sleep(5 * time.Second)
}

func boringWithChan(msg string, c chan string) {
	for i := 0; ; i++ {
		c <- fmt.Sprintf("%s %d", msg, i) // It waits for a receiver to be ready.
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func boring(msg string) {
	for i := 0; i < 3; i++ {
		fmt.Println(msg, i)
		time.Sleep(time.Second)
	}
}
