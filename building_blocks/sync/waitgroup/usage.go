// Package main illustrates the usage of waitgroup
package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup

func hello(wg *sync.WaitGroup, id int) {
	defer wg.Done()
	fmt.Printf("Hello from %v!\n", id)
}

const numGreeters = 6

func main() {
	wg.Add(1) //indicate that one goroutine is beginning
	go func() {
		// make waitgroup aware that a goroutine had finished
		defer wg.Done()

		fmt.Println("1st goroutine sleeping...")
		time.Sleep(1)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("2nd goroutine sleeping...")
		time.Sleep(2)
	}()

	wg.Add(numGreeters)
	for i := 0; i < numGreeters; i++ {
		go hello(&wg, i)
	}
	// Join Point back to main goroutine
	// it will block the main goroutine till all the goroutines
	// have indicated they have exited
	wg.Wait()
	fmt.Println("All goroutines are complete.")
}
