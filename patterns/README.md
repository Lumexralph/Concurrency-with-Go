# Concurrency Patterns

Ways to compose the Go concurrency primitives into patterns that will help keep your system scalable and maintainable.

## Confinement

It is the idea of ensuring information is ever available from one concurrent process. When this is achieved, a concurrent program is implicitly safe and no synchronization is needed.

They are of 2 types Ad hoc and lexical.

## The `for-select` Loop

It is very common in Go code, it is something like this;

        for { // Either loop infinitely or range over something
            select {
            // Do some work with channels
            case <- done:
                return
            default:
            // do something else
            }
        }

## Preventing Goroutine Leaks

Goroutines are not garbage collected by the runtime, so we don't want to leave them lying about our process. We have to clean them up.

How?

Establish a signal(channel) between parent goroutine and its children, that allows parent to signal cancellation to its children. Signal is read-only channel called `done` by convention. The parent goroutine passes this channel to the child goroutine and the closes the channel when it wants to cancel the child goroutine.

If a goroutine is responsible for creating a goroutine, it is also responsible for ensuring it can stop the goroutine.

## The `or-channel`

It is a pattern that creates a composite `done` channel through recursion and goroutines. Combining one or more `done` channels into a `done` channel that closes if any of its component channel closes.

## Error Handling

The most fundamental question when thinking about error handling is "Who should be responsible for handling the error?".

Your concurrent processes should send their errors to another part of the program that has complete information about the complete state of your program and can make more informed decision on what to do.

Errors should be considered first class citizen when constructing value to return from goroutines. If goroutines can produce errors, those errors should be tightly coupled with the result type, passed along through same line of communication.

## Pipelines

It is a very powerful tool to use when your program needs to process streams or batches of data. A pipeline is nothing more than a series of things that take data in, perform an operation on it and pass the data back out. Each of these operations are called `stages` of the pipeline.

Stages

Batch Processing - operate on a chunk of data all at once instead of one discrete value at a time.

Stream Processing - stage receives and returns one element at a time.

## Fan-Out, Fan-In

Fan-out is a term to describe the process of starting multiple goroutines to handle the input from the pipeline.

Fan-in is a term to describe the process of combining multiple results into one channel.

This pattern is used for the stages of a pipeline, when these two conditions apply;

* It doesn't rely on values that the stage has calculated before.

* It takes a long time to run.

## The or-done-channel

If the `done` channel of a goroutine we are reading another channel gets cancelled, we don't know for sure if the channel we are reading also gets cancelled. This patterns helps to handle situations like this.

## The tee-channel

In a situation you take in a stream of user commands and want to send them to something that executes them and also logs them for auditing, this pattern is a apt for this situation. It takes in a channel to read from and return two separate channels that will get same value.
