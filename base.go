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

func ChanToSensor[SensorEvent any](c chan SensorEvent) ISensor[SensorEvent] {
	return &BaseSensor[SensorEvent]{
		events: c,
	}
}

func NewSensor[SensorEvent any]() ISensor[SensorEvent] {
	return &BaseSensor[SensorEvent]{
		events: make(chan SensorEvent),
	}
}
