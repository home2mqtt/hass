package hass

import (
	"fmt"
	"strings"
)

type switch_impl struct {
	actuator baseActuator
	state    ISensor[BoolSensorEvent]
}

var _ ISwitch = &switch_impl{}

func NewSwitch(runtime IPubSubRuntime, config *Switch) ISwitch {
	state := NewSensor[BoolSensorEvent]()
	payloadon := fmt.Sprint(config.PayLoadOn)
	//payloadof := fmt.Sprint(config.PayLoadOff)
	runtime.Receive(config.StateTopic, func(topic string, payload []byte) {
		str := string(payload)
		on := strings.EqualFold(str, payloadon)
		state.Events() <- BoolSensorEvent{Value: on}
	})
	return &switch_impl{
		actuator: baseActuator{
			client: runtime,
			topic:  config.CommandTopic,
		},
		state: state,
	}
}

func (sw *switch_impl) State() ISensor[BoolSensorEvent] {
	return sw.state
}

func (sw *switch_impl) Toggle() {
	sw.actuator.send("{\"state\": \"TOGGLE\"}")
}

func (sw *switch_impl) Set(on bool) {
	var v string
	if on {
		v = "ON"
	} else {
		v = "OFF"
	}
	sw.actuator.send(fmt.Sprintf("{\"state\": \"%s\"}", v))
}
