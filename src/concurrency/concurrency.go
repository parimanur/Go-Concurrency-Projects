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

	time.Sleep(3 * time.Second)

	fmt.Println("Starting Fan In Pattern")
	runFanIn()
	fmt.Println("Finished Fan In Pattern\n\n")

	fmt.Println("Starting cleanup support using quit channel")
	runmMssagePassingWithCleanup()
	fmt.Println("Finished cleanup support using quit channel\n\n")

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

func runFanIn() {
	c := fanIn(boringChanGenerator("FooChan"), boringChanGenerator("BarChan"))
	for i := 0; i < 10; i++ {
		fmt.Println(<-c)
	}
	fmt.Println("You're both boring; I'm leaving.")
}

func fanIn(input1, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			c <- <-input1
		}
	}()
	go func() {
		for {
			c <- <-input2
		}
	}()
	return c
}

func fanInWithSelect(input1, input2 <-chan string) <-chan string {
	c := make(chan string)
	go func() {
		for {
			select {
			case s := <-input1:
				c <- s
			case s := <-input2:
				c <- s
			}
		}
	}()
	return c
}

func runmMssagePassingWithCleanup() {

	quit := make(chan string) // This channel is used for communication to quit and respond saying i am done from caller
	c := messagePassingWithCleanup("Joe", quit)

	for i := rand.Intn(10); i >= 0; i-- {
		fmt.Println(<-c)
	}

	quit <- "Bye!"                       // Quit signal main to caller
	fmt.Printf("Joe Says: %q\n", <-quit) // Caller responds after he finishes the cleanup
	//This is reached only after cleanup is done successfully
}

func messagePassingWithCleanup(msg string, quit chan string) <-chan string { // Returns receive-only (<-) channel of strings.
	c := make(chan string)

	go func() { // Launch the goroutine from inside the function. Function Literal.
		for i := 0; ; i++ {
			select {
			case c <- fmt.Sprintf("%s %d", msg, i):
				// Do Nothing.
			case <-quit:
				cleanup()
				quit <- "See you!"
				return
			}
		}
	}()

	return c // Return the channel to the caller.
}

func cleanup() {
	fmt.Println("Cleaned Up")
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
