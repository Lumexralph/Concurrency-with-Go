// Package main illustrated a fan-out, fan-in pattern using
// it to find prime numbers in a stream of data
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func randGenerator() interface{} {
	return rand.Intn(500_000_000)
}

func repeatFn(
	done <-chan interface{},
	fn func() interface{},
) <-chan interface{} {
	valueStream := make(chan interface{})
	go func() {
		defer close(valueStream)
		for {
			select {
			case <-done:
				return
			case valueStream <- fn():
			}
		}
	}()
	return valueStream
}

func take(
	done <-chan interface{},
	valueStream <-chan interface{},
	num int,
) <-chan interface{} {
	funcStream := make(chan interface{})
	go func() {
		defer close(funcStream)
		for i := num; i > 0 || i == -1; i-- {
			select {
			case <-done:
				return
			case funcStream <- <-valueStream:
			}
		}
	}()
	return funcStream
}

func toInt(
	done <-chan interface{},
	valueStream <-chan interface{},
) <-chan int {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for v := range valueStream {
			select {
			// cancellation/exit the goroutine
			case <-done:
				return
			// assert that the v type is an int and typecast it
			case intStream <- v.(int):
			}
		}
	}()
	return intStream
}

func oddFinder(
	done <-chan interface{},
	valueStream <-chan int,
) <-chan int {
	primeStream := make(chan int)
	go func() {
		defer close(primeStream)
		for v := range valueStream {

			select {
			// cancellation/exit the goroutine
			case <-done:
				return
			default:
				if v%2 == 1 {
					primeStream <- v
				}
			}
		}
	}()
	return primeStream
}

// fanIn joins the multiple stream of data into a single stream
func fanIn(
	done <-chan interface{},
	channels ...<-chan int,
) <-chan interface{} {
	// want to wait till all channels have been drained
	var wg sync.WaitGroup
	multiplexedStream := make(chan interface{})

	// read from the passed channel and put it into the
	// multiplexedStream
	multiplex := func(c <-chan int) {
		defer wg.Done()
		for i := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- i:
			}
		}
	}

	// select from all the channels amd wait for them
	wg.Add(len(channels))
	for _, c := range channels {
		go multiplex(c)
	}

	// wait for all the reads to complete
	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func main() {
	done := make(chan interface{})
	defer close(done)

	start := time.Now()
	randIntStream := toInt(done, repeatFn(done, randGenerator))

	// fan-out multiple stages of oddFinder
	numFinders := runtime.NumCPU()
	finders := make([]<-chan int, numFinders)
	for i := 0; i < numFinders; i++ {
		finders[i] = oddFinder(done, randIntStream)
	}

	fmt.Println("Primes:")
	for num := range take(done, fanIn(done, finders...), 100) {
		fmt.Printf("\t%d\n", num)
	}
	fmt.Printf("Search took: %v\n", time.Since(start))
}
