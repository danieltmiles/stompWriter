# stompWriter
A module to write to a queue using Golang's io.Writer interface

Usage:

```
stompWriter, err := New(hostname, port, username, password, queueName)
if err != nil {
	recover()
}

// pass stompWriter to any functions where Writer interface is supported

fmt.Fprint(stompWriter, "stomp message here")
```
