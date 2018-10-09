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

	time.Sleep(3 * time.Second)

	fmt.Println("Starting message passing using channel")
	runBoringWithChan()
	fmt.Println("Finished message passing using channel\n\n")

	time.Sleep(3 * time.Second)

	fmt.Println("Starting Generator Pattern")
	runBoringChanGenerator()
	fmt.Println("Finished Generator Pattern\n\n")

	time.Sleep(3 * time.Second)

	fmt.Println("Starting Sync Communication Pattern")
	runSyncCommunicationPattern()
	fmt.Println("Finished Sync Communication Pattern\n\n")

}

func runBoring() {

	// go statement runs the function as usual, but doesn't make the caller wait.
	go boring("boring go routine!!")
	// Running the below function will wait till the loop is completed
	//boring("boring main!")

	time.Sleep(3 * time.Second)
}
func runBoringWithChan() {

	c := make(chan string)
	go boringWithChan("chan c passed to go routine!", c)
	for i := 0; i < 5; i++ {
		fmt.Printf("Fetched from channel %q\n", <-c) // wait for a value to be sent.
	}
	fmt.Println("You're boring; I'm leaving.")
}

func runBoringChanGenerator() {
	c := boringChanGenerator("Generator Pattern!")
	for i := 0; i < 5; i++ {
		fmt.Printf("%q\n", <-c)
	}
}

func runSyncCommunicationPattern() {
	fooChan := boringChanGenerator("foo")
	barChan := boringChanGenerator("bar")

	for i := 0; i < 5; i++ {
		fmt.Println("Foo chan emits", <-fooChan)
		fmt.Println("bar chan emits", <-barChan) //This is blocked till fooChan emits
	}
}

func boring(msg string) {
	for i := 0; i < 3; i++ {
		fmt.Println(msg, i)
		time.Sleep(time.Second)
	}
}

func boringWithChan(msg string, c chan string) {
	for i := 0; ; i++ {
		c <- fmt.Sprintf("%s %d", msg, i) // It waits for a receiver to be ready.
		time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
	}
}

func boringChanGenerator(msg string) <-chan string {
	c := make(chan string)
	go func() {
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c
}
