# Concurrency Patterns

Ways to compose the Go concurrency primitives into patterns that will help keep your system scalable and maintainable.

## Confinement

It is the idea of ensuring information is ever available from one concurrent process. When this is achieved, a concurrent program is implicitly safe and no synchronization is needed.

They are of 2 types Ad hoc and lexical.
