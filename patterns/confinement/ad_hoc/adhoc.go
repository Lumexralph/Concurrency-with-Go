// Package main has the implementation of a ad hoc
// confinement concurrency pattern
package main

import "fmt"

var data = make([]int, 4)

// Even though the slice of data is accessible by the two functions
// it is available or accessed by the loopData goroutine
func loopData(handleData chan<- int) {
	// it is not the owner of the channel
	// I wonder why it has to close it
	defer close(handleData)
	for i := range data {
		handleData <- data[i]
	}
}

func main() {
	handleData := make(chan int)
	go loopData(handleData)

	for num := range handleData {
		fmt.Println(num)
	}
}
