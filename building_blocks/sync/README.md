# Sync Package

The package contains the primitives for low-level memory access synchronization.

## Various Primitives of the Package

1. WaitGroup - It is a great way to wait for a set of concurrent operations to complete when when you either don't care about the result of the concurrent operation or you have other means of collecting it, if not - use `channels` and select `statement` instead.
