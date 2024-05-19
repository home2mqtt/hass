package paho

import (
	"log"

	"github.com/balazsgrill/hass"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type msg struct {
	topic   string
	payload []byte
}

type PahoPubSub struct {
	mqtt.Client
	subscriptions map[string][]func(topic string, payload []byte)
	connstate     chan bool
	offlineMsgs   chan msg
}

var _ hass.IPubSubConnecterRuntime = &PahoPubSub{}

func New(options *mqtt.ClientOptions) hass.IPubSubConnecterRuntime {
	result := &PahoPubSub{
		subscriptions: make(map[string][]func(topic string, payload []byte)),
		connstate:     make(chan bool),
		offlineMsgs:   make(chan msg, 100),
	}
	options = options.SetOnConnectHandler(result.onConnect).SetConnectionLostHandler(result.onConnectionLost).SetAutoReconnect(true)

	result.Client = mqtt.NewClient(options)
	return result
}

func (ps *PahoPubSub) onConnectionLost(c mqtt.Client, err error) {
	if err != nil {
		log.Println(err)
	}
	ps.connstate <- false
}

func (ps *PahoPubSub) onConnect(client mqtt.Client) {
	for topic, callbacks := range ps.subscriptions {
		for _, callback := range callbacks {
			subscribe(ps.Client, topic, callback)
		}
	}
	for len(ps.offlineMsgs) > 0 {
		m := <-ps.offlineMsgs
		ps.internalsend(m.topic, m.payload, true)
	}
	ps.connstate <- true
}

func (ps *PahoPubSub) ConnectionState() chan bool {
	return ps.connstate
}

func (ps *PahoPubSub) Send(topic string, payload []byte) error {
	if ps.Client == nil || !ps.Client.IsConnected() {
		return nil
	}
	return ps.internalsend(topic, payload, false)
}

func (ps *PahoPubSub) internalsend(topic string, payload []byte, retained bool) error {
	token := ps.Client.Publish(topic, 0, retained, payload)
	token.Wait()
	return token.Error()
}

func (ps *PahoPubSub) SendRetained(topic string, payload []byte) error {
	if ps.Client == nil || !ps.Client.IsConnected() {
		ps.offlineMsgs <- msg{
			topic:   topic,
			payload: payload,
		}
		return nil
	}
	return ps.internalsend(topic, payload, true)
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
