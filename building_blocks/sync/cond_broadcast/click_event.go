// Package main illustrates creating a GUI application
// with a button. We want to register arbitrary number
// of functions when the button is clicked. We use the Cond
// Broadcast method to notify all registered handlers.
package main

import (
	"fmt"
	"sync"
)

// Button - is a type of GUI that contains the click Condition
type Button struct {
	Clicked *sync.Cond // need to send a click signal/event
}

var button = Button{Clicked: sync.NewCond(&sync.Mutex{})}

// subscribe allows to be able to register signals/event handlers
// each handler run its own goroutine.
func subscribe(c *sync.Cond, fn func()) {
	var runningGoroutine sync.WaitGroup
	runningGoroutine.Add(1)
	go func() {
		runningGoroutine.Done()
		// critical section
		c.L.Lock()
		defer c.L.Unlock()
		// suspend and block this goroutine until there is a
		// notification/signal
		c.Wait()
		fn()
		// end of critical section
	}()
	// block until all the goroutines have exited
	runningGoroutine.Wait()
}

var clickRegistered sync.WaitGroup

func main() {
	// run 3 goroutines, event handlers
	clickRegistered.Add(3)
	subscribe(button.Clicked, func() {
		defer clickRegistered.Done()
		fmt.Println("Loading the DOM of the page..")
	})

	subscribe(button.Clicked, func() {
		defer clickRegistered.Done()
		fmt.Println("Displaying alert box..")
	})

	subscribe(button.Clicked, func() {
		defer clickRegistered.Done()
		fmt.Println("Fetch all images..")
	})

	// main goroutine, sends a broadcast to all
	// goroutines (handlers) waiting to be notified
	// of the click event on the button that happened.
	button.Clicked.Broadcast() // click event has occurred.

	// block the main goroutine till the other goroutines have
	// finished or exited
	clickRegistered.Wait()
}
