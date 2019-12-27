# Go Concurrency

This is another study solely dedicated to understanding Concurrency and this will be learnt with Go, using the book Concurrency in Go by Katherine Cox-Buday as a guide.

I will be building small programs and making documentation of findings in any way possible.

## Concurrency and Parallelism

Concurrency is a property of the code, Parallelism is a property of the running program. We write a concurrent code and hope that it is run in parallel. A concurrent code running on one core is not running in parallel but sequentially in a fast way.

Parallelism is a property of the runtime of our code.

## Thoughts on Concurrency

Moore's law ("the number of components on an integrated circuit will double every 2 years") started losing it's impact till around the 2012, companies foresaw slowdown in rate of Moore's law and started looking for different ways in increasing computing power, this led to creating multicore processors, this led to solving problem in a simultaneous way (Parallelism) and this stem from Amdahl's law ("it entails, model the potential performance gains from implementing the solution to a problem in a parellel manner i.e gains are bounded by how much of the program must run in sequential manner").

Solving problems in a parallel manner led to Horizontal Scaling, having multiple instance of a program on different CPUs or machines.

## Sections

* Common Issues in Concurrency

* Building blocks of Go Concurrency
