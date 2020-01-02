// Package main is the implementation of the logic that monitors the
// health of a goroutine, to achieve this, it needs a reference to a
// function that can start the goroutine. It has a ward and a steward.
package main

import (
	"log"
	"os"
	"time"
)

type startGoroutineFn func(
	done <-chan interface{},
	pulseInterval time.Duration,
) (heartbeat <-chan interface{})

func newSteward(
	timeout time.Duration,
	startGoroutine startGoroutineFn,
) startGoroutineFn {
	return func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) <-chan interface{} {
		heartbeat := make(chan interface{})
		go func() {
			defer close(heartbeat)

			var wardDone chan interface{}
			var wardHearbeat <-chan interface{}
			// create a closure to encode a way to start goroutine
			// being monitored
			startWard := func() {
				// create a channel that gets passed to the ward
				// goroutine in case we want to signal it to be stopped
				wardDone = make(chan interface{})
				// start the goroutine to be monitored
				// to make the ward goroutine halt if either the
				// steward goroutine is halted or the steward halts
				// the ward goroutine, an or-channel is used to wrap
				// both done channels
				wardHearbeat = startGoroutine(or(wardDone, done), timeout/2)
			}
			startWard()
			pulse := time.Tick(pulseInterval)

		monitorLoop:
			for {
				timeoutSignal := time.After(timeout)
				// loop ensures the steward can send out pulses of its own
				for {
					select {
					case <-pulse:
						select {
						case heartbeat <- struct{}{}:
						default:
						}
					// receive the ward's pulse
					case <-wardHearbeat:
						continue monitorLoop
					// if we don't receive a pulse from the ward
					// within the set timeout, kill the ward goroutine
					// and restart a new ward goroutine
					case <-timeoutSignal:
						log.Println("steward: ward is unhealthy; restarting...")
						close(wardDone)
						startWard()
						continue monitorLoop
					case <-done:
						log.Println("steward: I am halting.")
						return
					}
				}
			}
		}()
		return heartbeat
	}
}

// takes variadic argument of channels and pack it into a slice
// var or func(channels ...<-chan interface{}) <-chan interface{}
func or(channels ...<-chan interface{}) <-chan interface{} {
	// base index
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	}

	orDone := make(chan interface{})
	go func() {
		defer close(orDone)

		switch len(channels) {
		case 2:
			select {
			case <-channels[0]:
			case <-channels[1]:
			}
		default:
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-channels[2]:
			// recursively create an or-channel from all the channels
			case <-or(append(channels[3:], orDone)...):
			}
		}
	}()
	return orDone
}

// the ward
func doWork(done <-chan interface{}, _ time.Duration) <-chan interface{} {
	log.Println("ward: Hello, I am not responsible for any work")

	go func() {
		<-done
		log.Println("ward: I am halting.")
	}()
	return nil
}

// the monitoring system
func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	monitorWithSteward := newSteward(4*time.Second, doWork)

	done := make(chan interface{})
	time.AfterFunc(9*time.Second, func() {
		log.Println("main: halting steward and ward.")
		close(done)
	})

	// start the  monitoring
	for range monitorWithSteward(done, 4*time.Second) {
	}
	log.Println("Done.")
}
