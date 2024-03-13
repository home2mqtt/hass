package hass

import (
	"fmt"
	"strconv"
)

type LightState struct {
	Brightness int    `json:"brightness,omitempty"`
	ColorMode  string `json:"color_mode,omitempty"`
	ColorTemp  int    `json:"color_temp,omitempty"`
	State      string `json:"state,omitempty"`
}

type light_impl struct {
	actuator   baseActuator
	state      light_state[bool]
	brightness *light_state[int]
	color_temp *light_state[int]
}

type light_state[Type any] struct {
	light  *baseActuator
	values chan Type
	key    string
}

func (state *light_state[Type]) send(command string) {
	state.light.send(fmt.Sprintf("{\"%s\":%s}", state.key, command))
}

var _ ILight = &light_impl{}

func NewLight(pubsub IPubSubRuntime, config Light) ILight {
	result := &light_impl{
		actuator: baseActuator{
			client: pubsub,
			topic:  config.CommandTopic,
		},
	}
	result.state = light_state[bool]{
		light:  &result.actuator,
		values: make(chan bool),
		key:    "state",
	}
	if config.Brightness {
		result.brightness = &light_state[int]{
			light:  &result.actuator,
			values: make(chan int),
			key:    "brightness",
		}
	}
	return result
}

func (d *light_impl) State() ISensor[bool] {
	return d.State()
}

func (d *light_impl) Toggle() {
	d.state.send("TOGGLE")
}
func (d *light_impl) Set(on bool) {
	var v string
	if on {
		v = "ON"
	} else {
		v = "OFF"
	}
	d.state.send(v)
}

func (d *light_impl) Brightness() IIntSettable {
	return d.brightness
}

func (d *light_impl) ColorTemp() IIntSettable {
	return d.color_temp
}

func (d *light_state[T]) SetValue(value int) {
	d.send(strconv.Itoa(value))
}

func (d *light_state[_]) Scale() int {
	return 254
}

func (d *light_state[T]) Values() chan T {
	return d.values
}
