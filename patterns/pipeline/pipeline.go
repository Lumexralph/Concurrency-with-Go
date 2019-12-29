// Package main illustrates a best practice pattern
// for creating pipelines
package main

import "fmt"

// generator - convert discrete set of values into a stream of data
// on a channel.
func generator(done <-chan interface{}, integers ...int) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for _, i := range integers {
			select {
			case <-done:
				return
			case intStream <- i:
			}
		}
	}()
	return intStream
}

func multiply(
	done <-chan interface{},
	intStream <-chan int,
	multiplier int,
) <-chan int {
	multipliedStream := make(chan int)
	go func() {
		defer close(multipliedStream)
		for i := range intStream {
			select {
			case <-done:
				return
			case multipliedStream <- i * multiplier:
			}
		}
	}()
	return multipliedStream
}

func add(
	done <-chan interface{},
	intStream <-chan int,
	additive int,
) <-chan int {
	addedStream := make(chan int)
	go func() {
		defer close(addedStream)
		for i := range intStream {
			select {
			case <-done:
				return // exit the goroutine(cancellation)
			case addedStream <- i + additive:
			}
		}
	}()
	return addedStream
}

func main() {
	// ensure our program exits cleanly
	// and never leaks goroutine
	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 2, 3, 4, 5, 6, 7)
	pipeline := multiply(done, add(done, multiply(done, intStream, 2), 10), 3)

	for v := range pipeline {
		fmt.Println(v)
	}
}
