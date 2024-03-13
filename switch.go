package hass

import (
	"fmt"
	"strings"
)

type switch_impl struct {
	actuator baseActuator
	state    ISensor[bool]
}

var _ ISwitch = &switch_impl{}

func NewSwitch(runtime IPubSubRuntime, config *Switch) ISwitch {
	payloadon := fmt.Sprint(config.PayLoadOn)
	payloadoff := fmt.Sprint(config.PayLoadOff)
	state := &BaseSensor[bool]{
		events: ParseSensorValue(runtime, config.StateTopic, config.ValueTemplate, func(s string) (bool, error) {
			if strings.EqualFold(s, payloadon) {
				return true, nil
			}
			if strings.EqualFold(s, payloadoff) {
				return false, nil
			}
			return false, fmt.Errorf("unknown boolean value: %s", s)
		}),
	}

	return &switch_impl{
		actuator: baseActuator{
			client: runtime,
			topic:  config.CommandTopic,
		},
		state: state,
	}
}

func (sw *switch_impl) State() ISensor[bool] {
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
