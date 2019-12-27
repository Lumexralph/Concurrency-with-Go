# Sync Package

The package contains the primitives for low-level memory access synchronization.

## Various Primitives of the Package

1. WaitGroup - It is a great way to wait for a set of concurrent operations to complete when when you either don't care about the result of the concurrent operation or you have other means of collecting it, if not - use `channels` and select `statement` instead.

2. Mutex (Mutual Exclusion) - It is a way to guard access to memory or critical section (areas where there is an important read or write to or from memory, area that requires an exclusive access to a shared resource) of the program. It is a way to give different aspect of the program to a shared resource one at a time. When using a Mutex, concurrent operations communicate by sharing memory, then memory access synchronization has to be done. It is about locking and unlocking access to the shared resource, it by chance you fail to unlock the Mutex, the program might run into a deadlock.

    When we are in a situation that not all the concurrent processes need to write or read, we can take advantage of the Read Write Mutex (sync.RWMutex).

3. Cond - A meeting or convergence point for goroutines waiting for or announcing the occurrence of an event(is an arbitrary signal between two or more goroutines that takes no information other than to say it has occurred).

    It is a way to make a goroutine efficiently sleep until it was signaled/event occurrence to wake and check its condition.

    Also, there is Broadcast method that provides a way to communicate with all the goroutines that are waiting the condition signal.

4. Once - It is a type that utilizes some `sync` primitives internally to ensure that only one call to `Do` ever calls the function passed in, even on different goroutines. The `sync.Once` only counts the number of times `Do` is called.

5. Pool - It is a concurrent-safe implementation of the object pool pattern (design patterns). It is a way to make available a fixed number or pool of things for use. It is used to constrain the creation of expensive things e.g database connections, a fixed number can be created but an indeterminate number of operations that request access to these things (pool).

    According to golang.org "Pool's purpose is to cache allocated but unused items for later reuse, relieving pressure on the garbage collector."

    To make more sense of this, if I have to use an object which is temporary and would be required by many concurrent processes, I can create a pool of these objects, pick from it to do what I want to do and return it back to the pool to avoid recreating such object which can be expensive to create in terms of time needed or might have impact on the memory allocated to our program to run.

    Another common situation where a Pool is used is for warming/making it ready a cache of pre-allocated objects for operations that must run as quickly as possible. To try front-load the time it takes to get a reference to another object.
