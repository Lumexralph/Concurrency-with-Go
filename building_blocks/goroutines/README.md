# Goroutines

Goroutines are one of the most basic units of organization in a Go program.

Every Go program has at least one goroutine: the main goroutine which is created and started when the process begins.

A goroutine is a function that is running concurrently alongside other code in a program.

Goroutines are not OS threads or threads managed by runtime (green threads), they are higher level of abstraction called Coroutines.

Coroutines are concurrent subroutines i.e functions, methods or closures in Go, that cannot be interrupted, instead they have multiple points that allows for suspension or reentry.
