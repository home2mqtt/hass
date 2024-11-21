package proxy

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/home2mqtt/hass"
)

var getStatePayload = []byte("{\"state\":\"\"}")

type iSubField interface {
	notify(string)
}

type jsonStateField struct {
	actuator  *baseActuator
	subfields map[string]iSubField
}

func (field *jsonStateField) NewIntField(key string, scale int) hass.IIntSettable {
	f := &intStateField{
		ISensor:   hass.NewSensor[int](),
		container: field,
		scale:     scale,
		key:       key,
	}
	field.subfields[key] = f
	return f
}

func (field *jsonStateField) NewOnOffField(key string) hass.IBoolField {
	f := &onOffStateField{
		ISensor:   hass.NewSensor[bool](),
		container: field,
		key:       key,
	}
	field.subfields[key] = f
	return f
}

func (field *jsonStateField) process(state map[string]interface{}) {
	for key, value := range state {
		if subfield, exists := field.subfields[key]; exists {
			subfield.notify(fmt.Sprint(value))
		}
	}
}

func (state *jsonStateField) send(key string, command string) {
	state.actuator.send(fmt.Sprintf("{\"%s\":\"%s\"}", key, command))
}

type intStateField struct {
	container *jsonStateField
	hass.ISensor[int]
	scale int
	key   string
}

func (d *intStateField) SetValue(value int) {
	d.container.send(d.key, strconv.Itoa(value))
}

func (field *intStateField) notify(value string) {
	v, err := strconv.Atoi(fmt.Sprint(value))
	if err != nil {
		field.SendEvent(v)
	}
}

func (field *intStateField) Scale() int {
	return field.scale
}

type onOffStateField struct {
	container *jsonStateField
	hass.ISensor[bool]
	key string
}

func (field *onOffStateField) Toggle() {
	field.container.send(field.key, "TOGGLE")
}

func (field *onOffStateField) notify(str string) {
	if strings.EqualFold("ON", str) {
		field.SendEvent(true)
	}
	if strings.EqualFold("OFF", str) {
		field.SendEvent(false)
	}
}

func (d *onOffStateField) SetValue(on bool) {
	var v string
	if on {
		v = "ON"
	} else {
		v = "OFF"
	}
	d.container.send(d.key, v)
}
