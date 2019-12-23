package main

import (
	"fmt"
	"sync"
)

// goroutines has a fork:join model
// fork - create a forked thread, goroutine (child)
// without the join spot
/*
func main() {
	go sayHello()
	go func() {
		fmt.Println("World!")
	}()
}

func sayHello() {
	fmt.Println("Hello!")
}
*/
var wg sync.WaitGroup

func main() {
	wg.Add(2)
	go sayHello()
	go func() {
		defer wg.Done()
		fmt.Println("World!")
	}()
	// blocks the main goroutine until other goroutines (child)
	wg.Wait() // Join point with the main goroutine
}

func sayHello() {
	defer wg.Done()
	fmt.Println("Hello!")
}
