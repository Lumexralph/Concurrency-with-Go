package main

import (
	"context"
	"fmt"
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

// Move concurrently fetches Things from fetch() and puts them in put(). It
// returns once fetch returns false (i.e. there are no more Things) and all
// Things have been put().
func Move(fetch Fetcher, put Putter) {
	// [TODO]

	go func() {
		// [TODO]
		t, ok := fetch()
		_, _ = t, ok // [Remove later]
		// [TODO]
	}()

	var t Thing // [Remove later]
	// [TODO]
	put(t)
	// [TODO]
}

// MaybeMove is exactly the same as Move except that it may return an error
// because fetch() and put() may return errors. If no errors occur then
// MaybeMove returns under the same conditions as Move(). If an error occurs
// then MaybeMove returns earlier even if there are more Things to fetch().
func MaybeMove(fetch MaybeFetcher, put MaybePutter) error {
	// [TODO]

	go func() {
		// [TODO]
		t, ok, err := fetch()
		_, _, _ = t, ok, err // [Remove later]
		// [TODO]
	}()

	var t Thing // [Remove later]
	// [TODO]
	err := put(t)
	_ = err // [Remove later]
	// [TODO]

	return nil // [Remove later]
}

// MoveCtx is exactly the same as Move except it honours the
// Context-cancellation channel returned by ctx.Done(). If ctx.Done() is closed
// early then MoveCtx returns ctx.Err() just as MaybeMove returns any errors
// that it encounters.
func MoveCtx(ctx context.Context, fetch Fetcher, put Putter) error {
	// [TODO]

	go func() {
		// [TODO]
		t, ok := fetch()
		_, _ = t, ok // [Remove later]
		// [TODO]
	}()

	var t Thing // [Remove later]
	// [TODO]
	put(t)
	// [TODO]

	return nil // [Remove later]
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
	fmt.Println(Fetcher())
}
