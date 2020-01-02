// Package main is the implementation of the logic that monitors the
// health of a goroutine, to achieve this, it needs a reference to a
// function that can start the goroutine. It has a ward and a steward.
package main

import (
	"fmt"
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

func doWorkFn(
	done <-chan interface{},
	intList ...int,
) (startGoroutineFn, <-chan interface{}) {
	intChanStream := make(chan (<-chan interface{}))
	intStream := bridge(done, intChanStream)

	// ward -
	doWork := func(
		done <-chan interface{},
		pulseInterval time.Duration,
	) <-chan interface{} {
		// channel to cumminicate on withing ward's goroutine
		intStream := make(chan interface{})
		hearbeat := make(chan interface{})

		go func() {
			defer close(intStream)

			select {
			// let bridge know the channel we'll be communicating on.
			case intChanStream <- intStream:
			case <-done:
				return
			}

			pulse := time.Tick(pulseInterval)

			for {
			valueLoop:
				for _, intVal := range intList {
					// simulate an unhealthy ward when negative value is seen
					if intVal < 0 {
						log.Printf("negative value: %d\n", intVal)
						return
					}

					for {
						select {
						case <-pulse:
							select {
							case hearbeat <- struct{}{}:
							default:
							}
						case intStream <- intVal:
							continue valueLoop
						case <-done:
							return
						}
					}
				}
			}
		}()
		return hearbeat
	}
	return doWork, intStream
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

func take(
	done <-chan interface{},
	valueStream <-chan interface{},
	num int,
) <-chan interface{} {
	funcStream := make(chan interface{})
	go func() {
		defer close(funcStream)
		for i := num; i > 0 || i == -1; i-- {
			select {
			case <-done:
				return
			case funcStream <- <-valueStream:
			}
		}
	}()
	return funcStream
}

// the monitoring system
func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ltime | log.LUTC)

	done := make(chan interface{})
	defer close(done)

	// create the ward
	doWork, intStream := doWorkFn(done, 3, 2, 1, 0, -1, 2, -3, 4, 3, 2, 1)
	// create the steward
	monitorWithSteward := newSteward(1*time.Millisecond, doWork)
	// start the ward and start monitoring
	monitorWithSteward(done, 1*time.Hour)

	for intVal := range take(done, intStream, 6) {
		fmt.Printf("main: received - %d\n", intVal)
	}
}
