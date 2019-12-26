# Sync Package

The package contains the primitives for low-level memory access synchronization.

## Various Primitives of the Package

1. WaitGroup - It is a great way to wait for a set of concurrent operations to complete when when you either don't care about the result of the concurrent operation or you have other means of collecting it, if not - use `channels` and select `statement` instead.

2. Mutex (Mutual Exclusion) - It is a way to guard access to memory or critical section (areas where there is an important read or write to or from memory, area that requires an exclusive access to a shared resource) of the program. It is a way to give different aspect of the program to a shared resource one at a time. When using a Mutex, concurrent operations communicate by sharing memory, then memory access synchronization has to be done. It is about locking and unlocking access to the shared resource, it by chance you fail to unlock the Mutex, the program might run into a deadlock.

When we are in a situation that not all the concurrent processes need to write or read, we can take advantage of the Read Write Mutex (sync.RWMutex).

3. Cond - A meeting or convergence point for goroutines waiting for or announcing the occurrence of an event(is an arbitrary signal between two or more goroutines that takes no information other than to say it has occurred).

It is a way to make a goroutine efficiently sleep until it was signaled/event occurrence to wake and check its condition.
