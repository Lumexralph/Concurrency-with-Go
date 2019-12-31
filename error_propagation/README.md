# Error Propagation

Error needs to relay a few pieces of critical information:

* What Happened - contains information about what led to the error.
* When and where it occurred - stack trace(how the call was initiated and where the error was instantiated), context it is running within (what machine the error occurred on, the UTC time the error was instantiated).
* A friendly user-facing message - about a line of text.
* How the user can get more information on the error.
