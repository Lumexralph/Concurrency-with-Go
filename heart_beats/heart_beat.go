// Package main illustrates a concurrent code, goroutine
// that exposes a heartbeat.
package main

import (
	"fmt"
	"time"
)

func doWork(
	done <-chan interface{},
	pulseInterval time.Duration,
) (<-chan interface{}, <-chan time.Time) {
	heartbeat := make(chan interface{}) // where heartbeats wii be sent
	results := make(chan time.Time)
	go func() {
		defer close(heartbeat)
		defer close(results)

		// at the every tick if the interval, there will something
		// to be read from these two channels
		pulse := time.Tick(pulseInterval)
		workGen := time.Tick(2 * pulseInterval)

		sendPulse := func() {
			select {
			case heartbeat <- struct{}{}:
			default: // we might not have a listener, to avoid block
			}
		}

		sendResult := func(r time.Time) {
			for {
				select {
				case <-done:
					return
				case <-pulse:
					sendPulse()
				case results <- r:
					return
				}
			}
		}

		for {
			select {
			case <-done:
				return
			case <-pulse:
				sendPulse()
			case r := <-workGen:
				sendResult(r)
			}
		}

	}()
	return heartbeat, results
}

// utilizing the heartbeat
func main() {
	done := make(chan interface{})
	// cancel the goroutines after 10 secs
	time.AfterFunc(10*time.Second, func() { close(done) })

	const timeout = 2 * time.Second
	heartbeat, results := doWork(done, timeout/2)

	for {
		select {
		case _, ok := <-heartbeat:
			if ok == false {
				return
			}
			fmt.Println("pulse")
		case r, ok := <-results:
			if ok == false {
				return
			}
			fmt.Printf("results %v\n", r.Second())
		// when the timeout elapses, end all process
		case <-time.After(timeout):
			return
		}
	}
}
