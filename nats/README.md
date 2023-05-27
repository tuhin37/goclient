# Nats client

## 1. Pakage

install the moduke

```go
go get github.com/tuhin37/goclient
```

import the nats package 

```go
import "github.com/tuhin37/goclient/nats"
```

## 

## 2. Create nast client

```go
natsClient, err := nats.NewNatsClient("127.0.0.1", "4222")
if err != nil {
    log.Fatal(err)
}
defer natsClient.Close()
```

## 

## 3. Define callback functions

`handleFoo()` callback function definition

```go
func handleFoo(payload []byte) {
    var data interface{}
    err := nats.DecodeMessage(payload, &data)
    if err != nil {
        log.Println("Error decoding message:", err)
        return
    }
    fmt.Println("Received message on subject 'foo':", data)
}
```

`handleBar()` callback function definition

```go
func handleBar(payload []byte) {
    var data interface{}
    err := nats.DecodeMessage(payload, &data)
    if err != nil {
        log.Println("Error decoding message:", err)
        return
    }
    fmt.Println("Received message on subject 'bar':", data)
}
```

## 

## 4. Subscribe & register callback functions

subscribe to `foo` topic and register `handleFoo()` callback function with it

```go
// handleFoo will be called when a new msg is available at "foo" topic
_, err = natsClient.Subscribe("foo", handleFoo)
if err != nil {
    log.Println("Error subscribing to subject 'foo':", err)
}
```

subscribe to `bar` topic and register `handleBar()` callback function with it

```go
// handleBar will be called when a new msg is available at "bar" topic
_, err = natsClient.Subscribe("bar", handleBar)
if err != nil {
    log.Println("Error subscribing to subject 'bar':", err)
}
```

    

## 5. Publish

publish to `foo` topic

```go
err = natsClient.PublishAny("foo", map[string]interface{}{"name": "tuhin", "age": 30})
if err != nil {
    log.Println("Error publishing message to subject 'foo':", err)
}
```

publish to `bar` topic 

```go
err = natsClient.PublishAny("bar", "yoyoyo")
if err != nil {
    log.Println("Error publishing message to subject 'bar':", err)
}
```

## 

## 6. Run forever

This keeps the program from exiting

```go
nats.Forever()
```

---



## 7. Complete code

```go
package main

import (
    "fmt"
    "log"

    "github.com/tuhin37/goclient/nats"
)

// Example subscription callback functions

func handleBar(payload []byte) {
    var data interface{}
    err := nats.DecodeMessage(payload, &data)
    if err != nil {
        log.Println("Error decoding message:", err)
        return
    }
    fmt.Println("Received message on subject 'bar':", data)
}

func handleFoo(payload []byte) {
    var data interface{}
    err := nats.DecodeMessage(payload, &data)
    if err != nil {
        log.Println("Error decoding message:", err)
        return
    }
    fmt.Println("Received message on subject 'foo':", data)
}

func main() {
    // Create a new NATS client
    natsClient, err := nats.NewNatsClient("127.0.0.1", "4222")
    if err != nil {
        log.Fatal(err)
    }
    defer natsClient.Close()

    // Subscribe to subjects
    _, err = natsClient.Subscribe("foo", handleFoo)
    if err != nil {
        log.Println("Error subscribing to subject 'foo':", err)
    }

    _, err = natsClient.Subscribe("bar", handleBar)
    if err != nil {
        log.Println("Error subscribing to subject 'bar':", err)
    }

    // Publish messages
    err = natsClient.PublishAny("foo", map[string]interface{}{"name": "tuhin", "age": 30})
    if err != nil {
        log.Println("Error publishing message to subject 'foo':", err)
    }

    err = natsClient.PublishAny("bar", map[int]interface{}{1: "tuhin", 2: 30})
    if err != nil {
        log.Println("Error publishing message to subject 'foo':", err)
    }

    // Wait for termination signal
    nats.Forever()
}
```

Output

```shell
Received message on subject 'bar': map[1:tuhin 2:30]
Received message on subject 'foo': map[age:30 name:tuhin]
```
