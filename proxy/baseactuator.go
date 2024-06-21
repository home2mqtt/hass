package proxy

import "github.com/home2mqtt/hass"

type baseActuator struct {
	client hass.IPubSubRuntime
	topic  string
}

func (b *baseActuator) init(context hass.IPubSubRuntime, topic string) {
	b.topic = topic
	b.client = context
}

func (s *baseActuator) send(action string) {
	s.client.Send(s.topic, []byte(action))
}
