// Package main illustrates cancellation of a goroutine
// blocked by on attempting to write value to a channel
package main

import (
	"fmt"
	"math/rand"
	"time"
)

func newRandStream(done <-chan interface{}) <-chan int {
	randStream := make(chan int)
	go func() {
		defer fmt.Println("newRandStream closure exited.")
		defer close(randStream)

		for {
			select {
			case randStream <- rand.Int():
			case <-done:
				return
			}
		}
	}()
	return randStream
}

func main() {
	done := make(chan interface{})
	randStream := newRandStream(done)
	fmt.Println("3 Random ints:")
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d: %d\n", i, <-randStream)
	}
	close(done)

	// simulate on-going work
	time.Sleep(1 * time.Second)
}
