package main

import (
	"fmt"
	"time"
)

// takes variadic slice of channels
// var or func(channels ...<-chan interface{}) <-chan interface{}
func or(channels ...<-chan interface{}) <-chan interface{} {
	// base index
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			// recursively create an or-channel from all the channels
			case <-or(append(channels[3:], orDone)...):
			}
		}
	}()
	return orDone
}

func sig(after time.Duration) <-chan interface{} {
	c := make(chan interface{})
	go func() {
		defer close(c)
		time.Sleep(after)
	}()
	return c
}

// wait for any of the channels
func main() {
	start := time.Now()
	// meeting point with main goroutine and other goroutines
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(3*time.Second),
		sig(1*time.Minute),
		sig(1*time.Hour),
	)
	fmt.Printf("done after %v\n", time.Since(start))
}
