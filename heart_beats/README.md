# Heartbeats

Heartbeats are a way for concurrent process to signal life to outside parties. It is a way to signal to its listeners that everything is well. and the silence is expected.

For any long-running goroutines, or goroutines that needs to be tested, this pattern is highly recommended.

There are two different types of heartbeats:

* Heartbeats that occur on a time interval.
* Heartbeats that occur at th beginning of a unit of work.