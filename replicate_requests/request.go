// Package main illustrates how to replicate a simulated requests
// over a certain number of handlers
package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func doWork(
	done <-chan interface{},
	id int,
	wg *sync.WaitGroup,
	result chan<- int,
) {
	started := time.Now()
	defer wg.Done()

	// simulate random load.
	simulatedLoadTime := time.Duration(1+rand.Intn(5)) * time.Second
	select {
	case <-done:
	case <-time.After(simulatedLoadTime):
	}

	select {
	case <-done:
	case result <- id:
	}

	took := time.Since(started)
	// display how long handler would take
	if took < simulatedLoadTime {
		took = simulatedLoadTime
	}
	fmt.Printf("%v took %v\n", id, took)
}

func main() {
	done := make(chan interface{})

	result := make(chan int)

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		// create multiple handlers to handle
		// that same request
		go doWork(done, i, &wg, result)
	}

	firstResponse := <-result
	close(done)
	wg.Wait()
	fmt.Printf("Received a response from #%v\n", firstResponse)
}
