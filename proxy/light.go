package proxy

import (
	"encoding/json"
	"log"

	"github.com/home2mqtt/hass"
)

type light_impl struct {
	actuator   baseActuator
	state      hass.IBoolField
	brightness hass.IIntSettable
	color_temp hass.IIntSettable
}

var _ hass.ILight = &light_impl{}

func NewLight(pubsub hass.IPubSubRuntime, config *hass.Light) hass.ILight {
	result := &light_impl{
		actuator: baseActuator{
			client: pubsub,
			topic:  config.CommandTopic,
		},
	}
	jsonStateField := &jsonStateField{
		actuator:  &result.actuator,
		subfields: map[string]iSubField{},
	}

	result.state = jsonStateField.NewOnOffField("state")
	if config.Brightness {
		result.brightness = jsonStateField.NewIntField("brightness", int(config.BrightnessScale))

	}
	if config.ColorMode {
		for _, colormode := range config.SupportedColorModes {
			switch colormode {
			case "color_temp":
				result.color_temp = jsonStateField.NewIntField("color_temp", 254)
			}
		}
	}

	pubsub.Receive(config.StateTopic, func(topic string, payload []byte) {
		var data map[string]interface{}
		err := json.Unmarshal(payload, &data)
		if err != nil {
			log.Println(err)
			return
		}
		jsonStateField.process(data)
	})
	pubsub.Send(config.StateTopic+"/get", getStatePayload)

	return result
}

func (d *light_impl) State() hass.ISensor[bool] {
	return d.state
}

func (d *light_impl) Toggle() {
	d.state.Toggle()
}
func (d *light_impl) Set(on bool) {
	d.state.SetValue(on)
}

func (d *light_impl) Brightness() hass.IIntSettable {
	return d.brightness
}

func (d *light_impl) ColorTemp() hass.IIntSettable {
	return d.color_temp
}
