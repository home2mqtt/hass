package hass

import (
	"fmt"
	"strconv"
)

type IPubSubRuntime interface {
	Send(topic string, payload []byte) error
	SendRetained(topic string, payload []byte) error
	Receive(topicpattern string, callback func(topic string, payload []byte))
}

type IConnecter interface {
	Connect() error
	Disconnect()
	ConnectionState() chan bool
}

type IPubSubConnecterRuntime interface {
	IPubSubRuntime
	IConnecter
}

func SendString(runtime IPubSubRuntime, topic string, payload string) error {
	return runtime.Send(topic, []byte(payload))
}

func SendFloat(runtime IPubSubRuntime, topic string, payload float64) error {
	return SendString(runtime, topic, fmt.Sprintf("%f", payload))
}

func ReceiveString(runtime IPubSubRuntime, topic string, callback func(topic string, payload string)) {
	runtime.Receive(topic, func(topic string, payload []byte) {
		callback(topic, string(payload))
	})
}

func ReceiveFloat(runtime IPubSubRuntime, topic string, callback func(topic string, payload float64, err error)) {
	ReceiveString(runtime, topic, func(topic, payload string) {
		value, err := strconv.ParseFloat(payload, 64)
		callback(topic, value, err)
	})
}
