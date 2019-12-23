// Package main aims to measure the amount of computer
// resources that was consumed to create goroutines
package main

import (
	"fmt"
	"runtime"
	"sync"
)

func memConsumed() uint64 {
	runtime.GC()
	// get the memory allocated to the program by the OS
	var s runtime.MemStats
	runtime.ReadMemStats(&s)
	return s.Sys
}

var c <-chan interface{}
var wg sync.WaitGroup

func noop() { wg.Done(); <-c }

const numGoroutines = 1e4

func main() {
	wg.Add(numGoroutines)
	before := memConsumed()

	for i := numGoroutines; i > 0; i-- {
		go noop()
	}
	wg.Wait() // Joint point to the main goroutine

	// we want to know the memory space that was consumed afterwards
	after := memConsumed()
	fmt.Printf("%.3fkb \n", float64(after-before)/numGoroutines/1000)
}
