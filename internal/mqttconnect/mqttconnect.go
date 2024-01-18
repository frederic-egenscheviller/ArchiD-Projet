package mqttconnect

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"os/signal"
)

type Client struct {
	client mqtt.Client
}

func NewClient(brokerAddress string, clientID string, defaultPublishHandler func(client mqtt.Client, message mqtt.Message)) (*Client, error) {
	opts := mqtt.NewClientOptions().AddBroker(brokerAddress)
	opts.SetClientID(clientID)
	if defaultPublishHandler != nil {
		opts.SetDefaultPublishHandler(defaultPublishHandler)
	}

	client := mqtt.NewClient(opts)

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &Client{client: client}, nil
}

func (m *Client) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	token := m.client.Publish(topic, qos, retained, payload)
	token.Wait()
	return token.Error()
}

func (m *Client) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	if token := m.client.Subscribe(topic, qos, callback); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (m *Client) Disconnect() {
	m.client.Disconnect(250)
	fmt.Println("Disconnected from MQTT broker")
}

func (m *Client) Unsubscribe(topic string) error {
	if token := m.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	fmt.Printf("Unsubscribed from topic %s\n", topic)
	return nil
}

func WaitForSignal() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	<-sigChan
}
