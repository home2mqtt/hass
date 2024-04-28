package hass

import (
	"encoding/json"
	"log"
)

type light_impl struct {
	actuator   baseActuator
	state      onOffStateField
	brightness *intStateField
	color_temp *intStateField
}

var _ ILight = &light_impl{}

func NewLight(pubsub IPubSubRuntime, config *Light) ILight {
	result := &light_impl{
		actuator: baseActuator{
			client: pubsub,
			topic:  config.CommandTopic,
		},
	}
	pubsub.Receive(config.StateTopic, func(topic string, payload []byte) {
		var data map[string]interface{}
		err := json.Unmarshal(payload, &data)
		if err != nil {
			log.Println(err)
			return
		}
		// TODO Check fields and propagate to state fields
		result.state.Process(data)
		if result.brightness != nil {
			result.brightness.Process(data)
		}
		if result.color_temp != nil {
			result.color_temp.Process(data)
		}
	})
	pubsub.Send(config.StateTopic+"/get", getStatePayload)

	result.state = onOffStateField{
		stateField: stateField[bool]{
			actuator: &result.actuator,
			values:   make(chan bool),
			key:      "state",
		},
	}
	if config.Brightness {
		result.brightness = &intStateField{
			stateField: stateField[int]{
				actuator: &result.actuator,
				values:   make(chan int),
				key:      "brightness",
			},
			scale: int(config.BrightnessScale),
		}
	}
	if config.ColorMode {
		for _, colormode := range config.SupportedColorModes {
			switch colormode {
			case "color_temp":
				result.color_temp = &intStateField{
					stateField: stateField[int]{
						actuator: &result.actuator,
						values:   make(chan int),
						key:      "color_temp",
					},
					scale: 254,
				}
			}
		}
	}
	return result
}

func (d *light_impl) State() ISensor[bool] {
	return &d.state
}

func (d *light_impl) Toggle() {
	d.state.send("TOGGLE")
}
func (d *light_impl) Set(on bool) {
	d.state.SetValue(on)
}

func (d *light_impl) Brightness() IIntSettable {
	return d.brightness
}

func (d *light_impl) ColorTemp() IIntSettable {
	return d.color_temp
}
