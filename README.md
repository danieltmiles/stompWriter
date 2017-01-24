# stompWriter
A module to write to a queue using Golang's io.Writer interface

Usage:

```
stompWriter, err := stompWriter.New(hostname, port, username, password, queueName)
if err != nil {
	//handle your error
}

// pass stompWriter to any functions where Writer interface is supported

fmt.Fprint(stompWriter, "stomp message here")
```
