// Package main has the implementation of a lexical
// confinement concurrency pattern, I favour this pattern
// over the ad hoc pattern
package main

import "fmt"

// encapsulate ownership, closing and writing to channel
// expose a read-only channel to any receiver
func chanOwner() <-chan int {
	// ownership, has lexical scope of this function
	// it also confines write access to the channel here
	results := make(chan int, 5)
	go func() {
		// closing
		defer close(results)
		for i := 0; i <= 5; i++ {
			results <- i
		}
	}()
	// implicit conversion of result bi-directional
	// channel to a read-only unidirectional channel
	return results
}

// receiving consumer handles the reading from channel
// and any blocking reason that might occur
func consumer(results <-chan int) {
	for num := range results {
		fmt.Printf("Received: %d\n", num)
	}
	fmt.Println("Done Receiving.")
}

func main() {
	results := chanOwner() // producer
	consumer(results)
}
