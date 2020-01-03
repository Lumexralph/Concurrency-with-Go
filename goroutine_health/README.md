# Healing UnHealthy Goroutines

It can be very easy for a goroutine to become stuck in a bad state from which it cannot recover without the help from an external source.
In a long-running process, it can be useful to create a mechanism that ensures your goroutines remain healthy and restarts them if they become unhealthy.

The process of restarting goroutines can be called `healing`. To heal goroutines, we can use the heartbeat pattern to check the liveliness of the goroutine being monitored.

The logic that monitors goroutine's health is called a `Steward`, while the goroutine that is being monitored is called the `ward`.
