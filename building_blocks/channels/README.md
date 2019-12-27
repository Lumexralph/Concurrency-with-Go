# Channels

Channels are one of the synchronization primitives in Go (memory access synchronization) - communicate to share memory.
Go's Concurrency is heavily influenced by Hoare's CSP report and this brought about channels. They are best used to communicate information between goroutines. Channels are typed.

A bidirectional (send and receive) channel

        var dataStream chan interface{}  // declare a type
        dataStream = make(chan interface{}) // instantiate the channel

## Unidirectional channels

Read or receive-only channels have the `<-` on the left-side of `chan`

        var dataStream <-chan interface{}  // declare a type
        dataStream = make(<-chan interface{}) // receive-only channel

Send or write-only channels have the `<-` on the right-side of `chan`

        var dataStream chan<- interface{}  // declare a type
        dataStream = make(chan<- interface{}) // send-only channel

Channels are blocking, just because a goroutine is scheduled doesn't guarantee it will run before the process exited.

Any goroutines that attempts to write to a channel that is full will wait till the channel has been emptied.

Also, any goroutine that attempts to read from a channel that is empty, it will wait till one item is placed on the channel. If the program is not structured properly using, it can cause deadlock.

## Buffered Channels

Channels that are given a certain capacity when they are instantiated. If the capacity is `n`, a goroutine can perform `n` number of writes without a read operation needed.

            var queue = make(chan int, 8) // can perform 8 writes

Buffered Channels are an in-memory FIFO queue for concurrent processes to communicate over.

## Nil Channel

Be very careful to avoid working with  a channel that has not been instantiated.

        var nilChan chan interface{}

Read and write to a nil channel leads to Fatal error, attempting to close a nil channel leads to panic!

Ensure channels are first initialized before any form of operation are done on it.

## Summary of Channel operations

Operation | Channel State | Result
----------|---------------|----------
Read    | nil              | Block
....... | Open and Not empty  | Value
....... | Open and empty  | Block
....... | Closed  | (default value), false
....... | Write Only | Compilation Error
....... | Open and Not empty  | Value
Write    | nil              | Block
....... | Open and Not Full  | Write Value
....... | Open and Full  | Block
....... | Closed  | Panic
....... | Receive Only | Compilation Error
Close    | nil              | Panic
....... | Open and Not empty  | Closes channel, reads succeed until channel is drained, then subsequent reads produce default value
....... | Open and empty  | Closes channel, then subsequent reads produce default value
....... | Closed  | Panic
....... | Receive Only | Compilation Error

src - Concurrency in Go book by Katherine Cox-Buday

## The select Statement

The `select` statement is the glue that binds channels together; it is how we are able to compose channels together in a program to form larger abstractions.

If there is no communication on the channels handled by a `select` statement, it blocks.

        select {
            case <-chan1:
                // do something
            case chan2<- "ball":
                // do something
            default:
                // do something if no channel is ready
        }

An empty `select` statement will block forever.

        select { }
