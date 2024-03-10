package hass

type BaseSensor[SensorEvent any] struct {
	events chan SensorEvent
}

func (s *BaseSensor[SensorEvent]) Close() error {
	close(s.events)
	return nil
}

func (s *BaseSensor[SensorEvent]) Events() chan SensorEvent {
	return s.events
}

func NewSensor[SensorEvent any]() ISensor[SensorEvent] {
	return &BaseSensor[SensorEvent]{
		events: make(chan SensorEvent),
	}
}

type baseActuator struct {
	client IPubSubRuntime
	topic  string
}

func (b *baseActuator) init(context IPubSubRuntime, topic string) {
	b.topic = topic
	b.client = context
}

func (s *baseActuator) send(action string) {
	s.client.Send(s.topic, []byte(action))
}
