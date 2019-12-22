// the code snippet is just to illustrate the issues
package issues

var store = [1]int{0}

// this goroutine is making a write operation
go func() {
	store[0]++
}()

// the if statement wants to readd the value stored
// since we can't guarantee the order of which part will
// be executed first, we have a race condition
if store[0] == 1 {
	fmt.Printf("The value is %d", store[0])
}
