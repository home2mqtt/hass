package hass

import (
	"fmt"
	"strconv"
	"strings"
)

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

type stateField[Type any] struct {
	actuator *baseActuator
	values   chan Type
	key      string
}

func (state *stateField[Type]) send(command string) {
	state.actuator.send(fmt.Sprintf("{\"%s\":\"%s\"}", state.key, command))
}

func (d *stateField[_]) Close() error {
	close(d.values)
	return nil
}

func (d *stateField[T]) Events() chan T {
	return d.values
}

type intStateField struct {
	stateField[int]
	scale int
}

func (d *intStateField) SetValue(value int) {
	d.send(strconv.Itoa(value))
}

func (field *intStateField) Process(state map[string]interface{}) {
	if value, exists := state[field.key]; exists {
		v, err := strconv.Atoi(fmt.Sprint(value))
		if err != nil {
			field.values <- v
		}
	}
}

func (field *intStateField) Scale() int {
	return field.scale
}

type onOffStateField struct {
	stateField[bool]
}

func (field *onOffStateField) Process(state map[string]interface{}) {
	if value, exists := state[field.key]; exists {
		str := fmt.Sprint(value)
		if strings.EqualFold("ON", str) {
			field.values <- true
		}
		if strings.EqualFold("OFF", str) {
			field.values <- false
		}
	}
}

func (d *onOffStateField) SetValue(on bool) {
	var v string
	if on {
		v = "ON"
	} else {
		v = "OFF"
	}
	d.send(v)
}
