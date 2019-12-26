// Package main illustrates a goroutine that is waiting for a signal
// and another goroutine that is sending the signals.
// We have have a queue of fixed length and some items on the queue,
// we want to enqueue items onto the queue as soon as possible when
// there is room. We want to get notified as soon as there is room
// in the queue to add new items.
package main

import (
	"fmt"
	"sync"
	"time"
)

// create a new condition
var c = sync.NewCond(&sync.Mutex{})

// instantiate a slice of length 0 and capacity 10
var queue = make([]interface{}, 0, 10)

func removeFromQueue(delay time.Duration) {
	time.Sleep(delay)
	// enter the critical section - accessing queue
	c.L.Lock()
	queue = queue[1:]
	fmt.Println("Removed from queue", len(queue))
	c.L.Unlock() // exiting the critical section
	c.Signal()   // notify any goroutine waiting on the condition
}

func main() {
	for i := 0; i <= 10; i++ {
		// enter the critical section i.e accessing shared
		// resource queue
		c.L.Lock()
		for len(queue) == 8 { // at the capacity we want
			// let's wait to dequeue, it will send a signal to continue
			// suspends the main goroutine until a signal on the
			// condition has been sent.
			c.Wait()
		}
		fmt.Println("Adding to the queue", len(queue))
		queue = append(queue, struct{}{})
		go removeFromQueue(1 * time.Second)
		c.L.Unlock() // exit the critical section
	}
	fmt.Printf("Final queue %v, len=%d, cap=%d\n", queue, len(queue), cap(queue))
}
