package main

import (
	"fmt"
	"time"
)

func doWork(
	done <-chan interface{},
	strings <-chan string,
) <-chan interface{} {
	terminated := make(chan interface{})
	go func() {
		defer fmt.Println("doWork exited.")
		defer close(terminated)

		for {
			select {
			case s := <-strings:
				fmt.Println(s)
			case <-done: // for cancelling this goroutine
				return
			}
		}
	}()
	return terminated
}

func main() {
	done := make(chan interface{})
	terminated := doWork(done, nil)

	// before we join doWork goroutine and main goroutine,
	// we have this goroutine to cancel doWork goroutine
	go func() {
		// cancel the doWork operation after 1 secs
		time.Sleep(1 * time.Second)
		fmt.Println("Cancelling doWork operation")
		close(done)
	}()

	// Join Point between goWork goroutine and the main goroutine
	<-terminated
	fmt.Println("Done")
}
