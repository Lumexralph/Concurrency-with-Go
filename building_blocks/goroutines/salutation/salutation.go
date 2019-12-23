package main

import (
	"fmt"
	"sync"
)

var wg sync.WaitGroup

func main() {
	for _, salutation := range []string{"hi", "good", "evening"} {
		wg.Add(1)
		// go func() {
		// 	defer wg.Done()

		// 	// the closure only binds to the last element
		// 	// this is because, the loop would have finished when
		// 	// goroutine starts executing, they all reference same
		// 	// memory space
		// 	fmt.Println(salutation)
		// }()

		// create a copy of every function parameter, different memory space
		go func(salutation string) {
			defer wg.Done()

			fmt.Println(salutation)
		}(salutation)
	}
	wg.Wait()
}
