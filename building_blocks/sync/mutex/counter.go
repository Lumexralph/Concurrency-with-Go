// Package main illustrated two goroutines that are
// attempting to increment or decrement a common value
// it uses a Mutex to synchronize access to memory
package main

import (
	"fmt"
	"sync"
)

var count int // shared resource by the goroutines
var lock sync.Mutex

var increment = func() {
	// critical area
	lock.Lock()
	defer lock.Unlock()
	count++
	fmt.Printf("Incrementing: %d\n", count)
}

var decrement = func() {
	// critical area
	lock.Lock()
	defer lock.Unlock()
	count--
	fmt.Printf("Decrementing: %d\n", count)
}

var arithmetic sync.WaitGroup

func main() {
	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			decrement()
		}()
	}

	for i := 0; i <= 5; i++ {
		arithmetic.Add(1)
		go func() {
			defer arithmetic.Done()
			increment()
		}()
	}
	// Join point between main and other goroutines
	// it blocks the main goroutine until all the
	// other goroutines have finished exiting
	arithmetic.Wait()
	fmt.Println("Calculation completed! final count", count)
}
