# Rate Limiting

It is a way of constraining the number of times some kind of resources is accessed to some finite number per unit of time. Resources can be anything like API connections, disk reads/writes, errors, network packets.

The reason this is mostly done is to prevent entire classes of attack vectors against our system. If you don't rate limit requests to your system, you cannot easily secure it. Also rate limiting a user requests can be advantageous to the application.

Most rate limiting is done by utilizing an algorithm called the [token bucket](https://en.wikipedia.org/wiki/Token_bucket). In production, there can be multiple layers of rate limiting.
