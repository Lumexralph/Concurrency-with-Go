# Issues in Concurrency

1. Race Condition -  It occurs when two or more operation must execute in the correct order, but the program has not been written so that the order is guaranteed or maintained.

2. Atomicity - issues with atomicity can be avoided through memory access synchronization

3. Deadlocks - A deadlocked program is one in which all the concurrent processes are waiting on one another. The program will never recover in this state without outside intervention.

4. Livelock - Programs that are actively performing concurrent operations, but these operations do nothing to move the state of the program forward.

5. Starvation - Is any situation where a concurrent process cannot get all the resources needed to perform work.
