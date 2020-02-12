package main

import (
	"context"
	"errors"
	"fmt"
	"time"
)

// A Thing is a dummy type that we are using for the exercise. It is "fetched"
// from a source and "put" in a destination.
type Thing interface{}

// A Fetcher function returns a Thing and a flag similar to the existence flag
// in a map—it will return true until there are no more Things and then return
// false.
type Fetcher func() (_ Thing, ok bool)

// A MaybeFetcher function is exactly the same as a Fetcher except that it may
// return an error.
type MaybeFetcher func() (_ Thing, ok bool, _ error)

// A Putter function accepts a Thing and stores it.
type Putter func(Thing)

// A MaybePutter function is exactly the same as a Putter except that it may
// return an error.
type MaybePutter func(Thing) error

// fetch function simulates returning Thing until there's nothing else
// it simulates a process that takes some time
var i int

type oldStore struct {
	bookNo        int
	bookInventory []Thing
}

type newStore struct {
	inventory []Thing
}

func (o *oldStore) fetch() (thing Thing, ok bool) {
	// simulate the delay
	time.Sleep(2 * time.Second)

	if o.bookNo > len(o.bookInventory)-1 {

		return thing, ok
	}

	time.Sleep(1 * time.Second)
	thing, ok = o.bookInventory[o.bookNo], true
	o.bookNo++
	return
}

// put take a value of type Thing and store it in a global map
// also simulates another long process
func (n *newStore) put(thing Thing) {
	time.Sleep(5 * time.Second)

	n.inventory = append(n.inventory, thing)
}

// Move concurrently fetches Things from fetch() and puts them in put(). It
// returns once fetch returns false (i.e. there are no more Things) and all
// Things have been put().
func Move(fetch Fetcher, put Putter) {
	ch := make(chan Thing)

	go func() {
		defer close(ch)

		for {
			t, ok := fetch()
			if !ok {
				break
			}
			ch <- t
		}
	}()

	// store the thing
	for thing := range ch {
		go put(thing)
	}
}

func (o *oldStore) fetchB() (thing Thing, ok bool, err error) {
	// simulate the delay
	time.Sleep(3 * time.Second)
	fmt.Println("fetchB: fetching...", o.bookNo)

	if o.bookNo == 5 {
		return thing, ok, errors.New("Fetch encountered network issues")
	}

	if o.bookNo > len(o.bookInventory)-1 {

		return thing, ok, err
	}

	thing, ok = o.bookInventory[o.bookNo], true
	o.bookNo++
	return
}

func (n *newStore) putB(thing Thing) error {
	time.Sleep(3 * time.Second)

	n.inventory = append(n.inventory, thing)

	if len(n.inventory) == 5 {
		return errors.New("store is filled up")
	}

	return nil
}

// MaybeMove is exactly the same as Move except that it may return an error
// because fetch() and put() may return errors. If no errors occur then
// MaybeMove returns under the same conditions as Move(). If an error occurs
// then MaybeMove returns earlier even if there are more Things to fetch().
func MaybeMove(fetch MaybeFetcher, put MaybePutter) error {
	ch := make(chan Thing)
	errCh := make(chan error, 2)
	quit := make(chan struct{})

	go func() {
		defer close(ch)

		for {
			select {
			// get a signal to stop the goroutine
			case <-quit:
				fmt.Println("Terminated from Parent Goroutine...")
				return
			default:
				t, ok, err := fetch()
				if err != nil {
					fmt.Println("Fetch encountered error")
					errCh <- err
					return
				}

				if !ok {
					return
				}

				ch <- t
			}
		}
	}()

	for thing := range ch {
		if err := put(thing); err != nil {
			fmt.Println("Error in adding to new store")
			close(quit) // signal to end the fetch goroutine
			errCh <- err
		}
	}

	close(errCh)

	return <-errCh
}

// MoveCtx is exactly the same as Move except it honours the
// Context-cancellation channel returned by ctx.Done(). If ctx.Done() is closed
// early then MoveCtx returns ctx.Err() just as MaybeMove returns any errors
// that it encounters.
func MoveCtx(ctx context.Context, fetch Fetcher, put Putter) error {
	ch := make(chan Thing)
	var fetchError error

	// I can avoid mixing anonymous function with goroutines
	go func() {
		defer close(ch)
		for {
			t, ok := fetch()
			if !ok {
				fmt.Println("===> I am done!")
				break
			}

			select {
			case <-ctx.Done():
				fetchError = ctx.Err()
				return
			case ch <- t:
			}
		}
	}()

	// store the thing
	for thing := range ch {
		put(thing)
	}
	return fetchError
}

// MoveLots functions like Move and also runs n concurrent go routines to fetch.
// It only returns once all of the go routines have returned and all the Things
// have been put().
func MoveLots(n int, fetch Fetcher, put Putter) {
	// [TODO]

	for i := 0; i < n; i++ {
		// [TODO]
		// Hint: sync.Waitgroup
		go func() {
			// [TODO]
			t, ok := fetch()
			_, _ = t, ok // [Remove later]
			// [TODO]
		}()
	}

	var t Thing // [Remove later]
	// [TODO]
	put(t)
	// [TODO]
}

// MaybeMoveLots combines the behaviour of all the other Move*() functions.
func MaybeMoveLots(ctx context.Context, n int, fetch MaybeFetcher, put MaybePutter) error {
	// [TODO]

	for i := 0; i < n; i++ {
		// [TODO]
		// Hint: errgroup.WithContext() instead of sync.Waitgroup
		go func() {
			// [TODO]
			t, ok, err := fetch()
			_, _, _ = t, ok, err // [Remove later]
			// [TODO]
		}()
	}

	var t Thing // [Remove later]
	// [TODO]
	err := put(t)
	_ = err // [Remove later]
	// [TODO]

	return nil // [Remove later]
}

func main() {
	old := &oldStore{
		bookNo: 0,
		bookInventory: []Thing{"concurrency with Go", "Go systems programming",
			"Isomorphic Go", "Go BluePrints", "Master Go",
			"Go Library Cookbook", "Algorithm and Data structures"},
	}
	new := &newStore{}

	// t := time.Now()
	// // we  can cancel the  whole process after a duration of 20 secs
	// ctx, cancelFn := context.WithCancel(context.Background())
	// time.AfterFunc(20*time.Second, cancelFn)

	// err := MoveCtx(ctx, old.fetch, new.put)
	// fmt.Println(new.inventory, "error: ", err)
	// fmt.Println("time Elapsed: ", time.Since(t))
	err := MaybeMove(old.fetchB, new.putB)
	fmt.Println("New Store: ", new.inventory)
	fmt.Println("MaybeMove: ", err)
}
