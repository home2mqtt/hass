package paho

import (
	"github.com/balazsgrill/hass"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type PahoPubSub struct {
	mqtt.Client
	subscriptions map[string][]func(topic string, payload []byte)
}

var _ hass.IPubSubConnecterRuntime = &PahoPubSub{}

func New(options *mqtt.ClientOptions) hass.IPubSubConnecterRuntime {
	result := &PahoPubSub{
		subscriptions: make(map[string][]func(topic string, payload []byte)),
	}
	options = options.SetOnConnectHandler(result.onConnect)
	result.Client = mqtt.NewClient(options)
	return result
}

func (ps *PahoPubSub) onConnect(client mqtt.Client) {
	for topic, callbacks := range ps.subscriptions {
		for _, callback := range callbacks {
			subscribe(ps.Client, topic, callback)
		}
	}
}

func (ps *PahoPubSub) Send(topic string, payload []byte) error {
	token := ps.Client.Publish(topic, 0, false, payload)
	token.Wait()
	return token.Error()
}

func subscribe(client mqtt.Client, topicpattern string, callback func(topic string, payload []byte)) {
	client.Subscribe(topicpattern, 0, func(c mqtt.Client, m mqtt.Message) {
		callback(m.Topic(), m.Payload())
	}).Wait()
}

func (ps *PahoPubSub) Receive(topicpattern string, callback func(topic string, payload []byte)) {
	list, ok := ps.subscriptions[topicpattern]
	if ok {
		list = append(list, callback)
	} else {
		list = []func(topic string, payload []byte){callback}
	}
	ps.subscriptions[topicpattern] = list
	if ps.Client.IsConnected() {
		subscribe(ps.Client, topicpattern, callback)
	}
}

func (ps *PahoPubSub) Connect() error {
	token := ps.Client.Connect()
	token.Wait()
	return token.Error()
}

func (ps *PahoPubSub) Disconnect() {
	ps.Client.Disconnect(100)
}
