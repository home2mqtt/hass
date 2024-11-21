package hass

type BaseSensor[SensorEvent any] struct {
	handler func(SensorEvent)
}

func (s *BaseSensor[SensorEvent]) ReceiveEvent(handler func(SensorEvent)) {
	s.handler = handler
}

func (s *BaseSensor[SensorEvent]) SendEvent(event SensorEvent) {
	if s.handler != nil {
		s.handler(event)
	}
}

func ChanToSensor[SensorEvent any](c chan SensorEvent) ISensor[SensorEvent] {
	result := &BaseSensor[SensorEvent]{}
	go func() {
		for {
			s := <-c
			result.SendEvent(s)
		}
	}()
	return result
}

func NewSensor[SensorEvent any]() ISensor[SensorEvent] {
	return &BaseSensor[SensorEvent]{}
}
