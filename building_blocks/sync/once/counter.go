package main

import (
	"fmt"
	"sync"
)

var count int

func increment() { count++ }
func decrement() { count-- }

var once sync.Once

func main() {
	once.Do(increment)
	once.Do(decrement)

	fmt.Printf("Count: %d\n", count) // count is 1
}
