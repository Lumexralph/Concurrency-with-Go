# Replicated Requests

This pattern is a good idea for some applications that receiving a response as quickly as possible is a top priority. You can replicate the request to multiple handlers(processes, goroutines or servers), and one of them will run faster than the other ones, the result can immediately be returned.

The downside to this approach is that you'll have to utilize resources to keep multiple copies of the handlers running, whatever resources the  handlers are using to service the requests need to be replicated as well.
