package main

import (
	"fmt"
	"sync"
	"time"
)

var wg sync.WaitGroup
var sharedLock sync.Mutex

const runtime = 1 * time.Second

var greedyWorker = func() {
	defer wg.Done()

	var count int
	for begin := time.Now(); time.Since(begin) <= runtime; {
		sharedLock.Lock()
		time.Sleep(3 * time.Nanosecond)
		sharedLock.Unlock()
		count++
	}
	fmt.Printf("Greedy worker was able to execute %v work loops \n", count)
}

var politeWorker = func() {
	wg.Done()

	var count int
	for begin := time.Now(); time.Since(begin) <= runtime; {
		sharedLock.Lock()
		time.Sleep(1 * time.Nanosecond)
		sharedLock.Unlock()

		sharedLock.Lock()
		time.Sleep(1 * time.Nanosecond)
		sharedLock.Unlock()

		sharedLock.Lock()
		time.Sleep(1 * time.Nanosecond)
		sharedLock.Unlock()

		count++
	}
	fmt.Printf("Polite worker was able to execute %v work loops \n", count)
}

func main() {
	wg.Add(2)
	go greedyWorker() // take so much of the time leaving little for the polite worker
	go politeWorker()

	wg.Wait()
}