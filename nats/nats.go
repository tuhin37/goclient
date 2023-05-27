package nats

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
)

// NatsClient represents a client for NATS messaging.
type NatsClient struct {
	nc *nats.Conn
}

// NewNatsClient creates a new NATS client and establishes a connection to the server.
func NewNatsClient(host, port string) (*NatsClient, error) {
	nc, err := nats.Connect(fmt.Sprintf("nats://%s:%s", host, port))
	if err != nil {
		return nil, err
	}

	return &NatsClient{
		nc: nc,
	}, nil
}

// Publish publishes a message to the specified subject.
func (client *NatsClient) Publish(subject string, data interface{}) error {
	payload, err := encodeData(data)
	if err != nil {
		return err
	}

	return client.nc.Publish(subject, payload)
}

// PublishAny publishes a message to the specified subject, supporting any kind of data.
func (client *NatsClient) PublishAny(subject string, data interface{}) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return client.nc.Publish(subject, payload)
}

// Subscribe registers a callback function to handle messages received on the specified subject.
func (client *NatsClient) Subscribe(subject string, callback func([]byte)) (*nats.Subscription, error) {
	return client.nc.Subscribe(subject, func(msg *nats.Msg) {
		callback(msg.Data)
	})
}

// DecodeMessage decodes the received data into the provided object.
func DecodeMessage(data []byte, obj interface{}) error {
	return json.Unmarshal(data, obj)
}

// Close closes the NATS connection.
func (client *NatsClient) Close() {
	client.nc.Close()
}

// Helper function to encode data
func encodeData(data interface{}) ([]byte, error) {
	switch t := data.(type) {
	case string:
		return []byte(t), nil
	default:
		return json.Marshal(data)
	}
}

func Forever() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
