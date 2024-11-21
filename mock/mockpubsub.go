package mock

import (
	"github.com/home2mqtt/hass"
)

type MockPubSub struct {
	callbacks map[string]func(string, []byte)
}

func NewMockPubSub() *MockPubSub {
	return &MockPubSub{
		callbacks: make(map[string]func(string, []byte)),
	}
}

var _ hass.IPubSubRuntime = &MockPubSub{}

func (tc *MockPubSub) Send(topic string, command []byte) error {
	if tc.callbacks == nil {
		return nil
	}
	if c, ok := tc.callbacks[topic]; ok {
		c(topic, command)
	}
	return nil
}

func (tc *MockPubSub) SendRetained(topic string, command []byte) error {
	return tc.Send(topic, command)
}

func (tc *MockPubSub) Receive(topic string, callback func(string, []byte)) {
	if tc.callbacks == nil {
		tc.callbacks = make(map[string]func(string, []byte))
	}
	tc.callbacks[topic] = callback
}
