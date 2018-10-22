/*

Q: What does Google Search do?

A: Given a query, return a page of Search Results (and some ads).

Q: How do we get the Search Results?

A: Send the query to Web Search, Image Search, YouTube, Maps, News,etc., then mix the Results.

*/
package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	Web   = fakeSearch("web1")
	Image = fakeSearch("image1")
	Video = fakeSearch("video1")
	//Below are the replicated search servers.
	Web2   = fakeSearch("web2")
	Image2 = fakeSearch("image2")
	Video2 = fakeSearch("video2")
)

type (
	Result string //result returned by the individual search modules
	Search func(query string) Result
)

func main() {
	runVer1()
	runVer2()
	runVer3()
}

func runVer1() {
	rand.Seed(time.Now().UnixNano())

	start := time.Now()
	Results := google_ver1("golangv1")
	elasped := time.Since(start)

	fmt.Println(Results)
	fmt.Println("Elapsed duration for synchronised execution of Search: ", elasped, "\n\n")
}

func runVer2() {
	rand.Seed(time.Now().UnixNano())

	start := time.Now()
	Results := google_ver2("golangv2")
	elasped := time.Since(start)

	fmt.Println(Results)
	fmt.Println("Elapsed duration for concurrent but  not resilient execution of Search: ", elasped, "\n\n")
}

func runVer3() {
	rand.Seed(time.Now().UnixNano())

	start := time.Now()
	Results := google_ver3("golangv3")
	elasped := time.Since(start)

	fmt.Println(Results)
	fmt.Println("Elapsed duration for fast concurrent resilient and robust of Search: ", elasped, "\n\n")
}

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result(fmt.Sprintf("%s Result for %q\n", kind, query))
	}
}

func google_ver1(query string) (Results []Result) {
	Results = append(Results, Web(query))
	Results = append(Results, Image(query))
	Results = append(Results, Video(query))

	return Results
}

func google_ver2(query string) (Results []Result) {
	c := make(chan Result)
	//Run the Web, Image, and Video Searches concurrently, and wait for all Results.
	//No locks. No condition variables. No callbacks.

	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	for i := 0; i < 3; i++ {
		Result := <-c
		Results = append(Results, Result)
	}
	return
}

func google_search_with_timeout(query string) (Results []Result) {
	c := make(chan Result)
	go func() { c <- Web(query) }()
	go func() { c <- Image(query) }()
	go func() { c <- Video(query) }()

	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case Result := <-c:
			Results = append(Results, Result)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return
}

/* Avoid timeout
Q: How do we avoid discarding Results from slow servers?

A: Replicate the servers. Send requests to multiple replicas, and use the first response.*/

func First(query string, replicas ...Search) Result {
	c := make(chan Result)
	//inline function to call the module search web/video ...
	searchReplica := func(i int) { c <- replicas[i](query) }
	for i := range replicas {
		go searchReplica(i)
	}
	//As soon as any of the replica modules return first return that value
	return <-c
}

func google_ver3(query string) (results []Result) {

	c := make(chan Result)
	go func() { c <- First(query, Web, Web2) }()
	go func() { c <- First(query, Image, Image2) }()
	go func() { c <- First(query, Video, Video2) }()
	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case Result := <-c:
			results = append(results, Result)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return
}
