// Package main illustrated goroutines that take the place
// of ownership of a channeland another that handles reading from it.
package main

import "fmt"

var chanOwner = func() <-chan int {
	// encapsulates all the operation of instantiating
	// closing and writes to a channel in a read channel
	resultStream := make(chan int)
	go func() {
		// handle close of the channel
		defer close(resultStream)

		for i := 0; i <= 5; i++ {
			// handle write to the channel
			resultStream <- i
		}
	}()

	// the bidirectional channel is implicitly converted
	// to a read-only channel
	return resultStream
}

func main() {
	// consumer goroutine or reader goroutine
	// handles when a channel is closed and any
	// blocking reason
	resultStream := chanOwner()
	for result := range resultStream {
		fmt.Printf("Received: %d\n", result)
	}
	fmt.Println("Done receiving.")
}
