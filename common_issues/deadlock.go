// Note:  the code snippet is just for illustration reasons
// it might not actually work
package main

import (
	"fmt"
	"sync"
	"time"
)

type shop struct {
	mu    sync.Mutex
	store [1]int
}

func main() {
	var wg sync.WaitGroup
	doubleQuantity := func(s *shop) {
		defer wg.Done()
		// critical section 1
		s.mu.Lock()
		s.store[0] = s.store[0] * 2
		defer s.mu.Unlock()

		time.Sleep(2 * time.Second)
		// critical section 2
		s.mu.Lock()
		fmt.Printf("The value is %d", s.store[0])
		defer s.mu.Unlock()
	}

	var shopRite shop
	wg.Add(2)
	go doubleQuantity(&shopRite)
	go doubleQuantity(&shopRite)
	wg.Wait()

}
