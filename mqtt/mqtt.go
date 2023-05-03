package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MqttClient struct {
	client mqtt.Client
}

func NewMqttClient(host string, port string, username string, password string) (MqttClient, error) {
	m := MqttClient{}
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%s", host, port))
	if username != "" && password != "" {
		opts.SetUsername(username)
		opts.SetPassword(password)
	}
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	m.client = mqtt.NewClient(opts)

	token := m.client.Connect()
	if token.Wait() && token.Error() != nil {
		return MqttClient{}, token.Error()
	}

	return m, nil
}

func (m *MqttClient) Subscribe(topic string) {
	token := m.client.Subscribe(topic, 1, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic %s", topic)
}

func (m *MqttClient) Disconnect(quiesce uint) {
	m.client.Disconnect(quiesce)
	fmt.Println("Disconnected")
}

func (m *MqttClient) Publish(topic string, qos byte, retained bool, payload interface{}) {
	token := m.client.Publish(topic, qos, retained, payload)
	token.Wait()

}

// Helper functions
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received message: %s from topic: %s\n", msg.Payload(), msg.Topic())
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v", err)
}

// this will get shared around the app
var Mqtt MqttClient

// func init() {
// 	util.LoadENV()
// 	Mqtt.Initialize()
// 	Mqtt.Subscribe(os.Getenv("MQTT_TOPIC_MSG_DOWN"))
// }
