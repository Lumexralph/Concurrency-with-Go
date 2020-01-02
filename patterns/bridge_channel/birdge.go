// Package main is the implementation of a bridge-channel pattern
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

// it helps destructuring a channel of channels into a single channel
func bridge(
	done <-chan interface{},
	chanStream <-chan <-chan interface{},
) <-chan interface{} {
	// single channel to return all values
	// from the stream of channels
	valStream := make(chan interface{})
	go func() {
		defer close(valStream)
		// pull values(chan) of the stream of channels
		for {
			var stream <-chan interface{}
			select {
			case maybeStream, ok := <-chanStream:
				if ok == false {
					return
				}
				stream = maybeStream
			case <-done:
				return
			}

			for val := range orDone(done, stream) {
				select {
				case valStream <- val:
				case <-done:
				}
			}
		}
	}()
	return valStream
}

// example to utilize the use of bridge, it returns a sequence
// of channels.
func genVals() <-chan <-chan interface{} {
	// create a chan that expects a read-only channel type
	chanStream := make(chan (<-chan interface{}))
	go func() {
		defer close(chanStream)

		for i := 0; i < 10; i++ {
			stream := make(chan interface{}, 1)
			stream <- i
			close(stream)
			chanStream <- stream
		}
	}()
	return chanStream
}

func main() {
	done := make(chan interface{})
	defer close(done)

	for v := range bridge(done, genVals()) {
		fmt.Printf("%d ", v)
	}
	fmt.Println("")
}
