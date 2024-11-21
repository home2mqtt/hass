package proxy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/home2mqtt/hass"
)

type boolField struct {
	baseActuator
	hass.ISensor[bool]
	on  string
	off string
}

// SetValue implements hass.IBoolField.
func (b *boolField) SetValue(value bool) {
	if value {
		b.send(b.on)
	} else {
		b.send(b.off)
	}
}

// Toggle implements hass.IBoolField.
func (b *boolField) Toggle() {
	b.send("TOGGLE")
}

func (b *boolField) SetOnValue(on string) {
	b.on = on
}

func (b *boolField) SetOffValue(off string) {
	b.off = off
}

func NewBoolField(runtime hass.IPubSubRuntime, topic string, cmdtopic string) hass.IBoolField {
	f := &boolField{
		baseActuator: baseActuator{
			client: runtime,
			topic:  cmdtopic,
		},
		ISensor: hass.NewSensor[bool](),
		on:      "ON",
		off:     "OFF",
	}
	if topic != "" {
		runtime.Receive(topic, func(topic string, payload []byte) {
			value := string(payload)
			onoff := strings.EqualFold(value, "on")
			f.SendEvent(onoff)
		})
	}
	return f
}

type enumField struct {
	baseActuator
	hass.ISensor[string]
	values []string
}

// List implements hass.IEnumField.
func (e *enumField) List() []string {
	return e.values
}

// SetValue implements hass.IEnumField.
func (e *enumField) SetValue(value string) {
	e.send(value)
}

func NewStringSensor(runtime hass.IPubSubRuntime, valuetemplate string, topic string) hass.ISensor[string] {
	return hass.ChanToSensor(ParseSensorValue(runtime, topic, valuetemplate, func(s string) (string, error) {
		return s, nil
	}))
}

func NewEnumField(runtime hass.IPubSubRuntime, topic string, cmdtopic string, values []string) hass.IEnumField {
	f := &enumField{
		baseActuator: baseActuator{
			client: runtime,
			topic:  cmdtopic,
		},
		ISensor: NewStringSensor(runtime, "", topic),
		values:  values,
	}
	return f
}

type floatField struct {
	baseActuator
	hass.ISensor[float64]
}

// SetValue implements hass.IField.
func (f *floatField) SetValue(value float64) {
	f.send(fmt.Sprint(value))
}

func NewFloatSensor(runtime hass.IPubSubRuntime, valuetemplate string, topic string) hass.ISensor[float64] {
	if topic == "" {
		return nil
	}
	return hass.ChanToSensor(ParseSensorValue(runtime, topic, valuetemplate, func(s string) (float64, error) {
		v, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return 0, err
		}
		return v, nil
	}))
}

func NewFloatField(runtime hass.IPubSubRuntime, valuetemplate string, topic string, cmdtopic string) hass.IField[float64] {
	f := &floatField{
		baseActuator: baseActuator{
			client: runtime,
			topic:  cmdtopic,
		},
		ISensor: NewFloatSensor(runtime, valuetemplate, topic),
	}
	return f
}
