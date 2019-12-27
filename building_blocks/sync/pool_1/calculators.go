// Package main illustrates the use of sync.Pool,
// Pool's primary interface is its Get method. When called,
// Get will first check whether there are any available instances
// in the pool to return to the caller, and if not, call its New
// member vairable to create one.
package main

import (
	"fmt"
	"sync"
)

var numCalcCreated int

func calcPool() *sync.Pool {
	p := &sync.Pool{
		// there must always be a New member field
		// that returns a pointer to the created object
		New: func() interface{} {
			numCalcCreated++
			mem := make([]byte, 1024)
			return &mem
		},
	}
	return p
}

func main() {
	// seed the pool with 4KB of []bytes
	calcPool().Put(calcPool().New())
	calcPool().Put(calcPool().New())
	calcPool().Put(calcPool().New())
	calcPool().Put(calcPool().New())

	const numWorkers = 1024 * 1024
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for i := numWorkers; i > 0; i-- {
		go func() {
			defer wg.Done()
			// do something or some operation with this memory
			// assert that the type returned is a pointer to
			// a slice of bytes and convert the empty interface
			mem := calcPool().Get().(*[]byte)
			// when you're done using this memory,
			// return it back to the pool to be reused
			defer calcPool().Put(mem)
		}()
	}
	// block the main goroutine till other goroutines have finished
	wg.Wait()
	fmt.Printf("%d calculators were created\n", numCalcCreated)
}
