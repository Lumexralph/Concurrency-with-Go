// Package main is the implementation of the or-done-channel pattern
package main

import "fmt"

func orDone(done, c <-chan interface{}) <-chan interface{} {
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		for {
			select {
			case <-done:
				return
			case v, ok := <-c:
				// if channel has been closed
				if ok == false {
					return
				}
				// continue reading value from the channel
				select {
				case valStream <- v:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

func main() {
	// we can then do something with the channel we are reading from
	done := make(chan interface{})
	aChan := make(chan interface{})
	for val := range orDone(done, aChan) {
		fmt.Println(val)
	}

}
