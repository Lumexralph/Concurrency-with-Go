package issues

import "sync"
// to ensure that a critical section (parts of the program)
// that needs access to the memory have exclusive access to
// the shared memory at different point in time, a memory access
// synchronization can be done

// Note: The code is for illustration, might not be working
// and This is not idiomatic Go because
// - we have not avoided the race condition
// - Also locking has some performance perks, whenever
// we lock, our program pauses for a period of time.

// Before memory access synchronization
var store = [1]int{0}

// critical section 1
go func() {
	store[0]++
}

// critical section 2
if store[0] == 1 {
	fmt.Printf("The value is %d", store[0])
}

// After memory access synchronization
var memoryAccess sync.Mutex
var store = [1]int{0}

// critical section 1
go func() {
	memoryAccess.Lock()
	store[0]++
	memoryAccess.Unlock()
}

// critical section 2
memoryAccess.Lock()
if store[0] == 1 {
	fmt.Printf("The value is %d", store[0])
}
memoryAccess.Unlock()
