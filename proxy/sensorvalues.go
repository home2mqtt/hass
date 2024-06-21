package proxy

import (
	"encoding/json"
	"log"

	"github.com/home2mqtt/hass"
	"github.com/noirbizarre/gonja"
	"github.com/noirbizarre/gonja/exec"
)

func GenericSensor(context hass.IPubSubRuntime, sensor *hass.Sensor) hass.ISensor[string] {
	values := ParseSensorValue(context, sensor.Topic, sensor.ValueTemplate, func(s string) (string, error) {
		return s, nil
	})
	result := hass.ChanToSensor(values)

	return result
}

func ParseSensorValue[T any](context hass.IPubSubRuntime, Topic string, ValueTemplate string, parse func(string) (T, error)) chan T {
	result := make(chan T)

	var tpl *exec.Template
	if ValueTemplate != "" {
		var err error
		tpl, err = gonja.FromString(ValueTemplate)
		if err != nil {
			log.Println(err)
		}
	}
	context.Receive(Topic, func(topic string, payload []byte) {
		var value string
		if tpl == nil {
			value = string(payload)
		} else {
			var data map[string]interface{}
			err := json.Unmarshal(payload, &data)
			if err != nil {
				log.Println(err)
				return
			}
			context := map[string]interface{}{"value_json": data}
			value, err = tpl.Execute(context)
			if err != nil {
				log.Println(err)
				return
			}
		}
		e, err := parse(value)
		if err != nil {
			log.Println(err)
			return
		}
		select {
		case result <- e:
			break
		default:
			// Drop data if no readers are waiting
			break
		}
	})
	context.Send(Topic+"/get", getStatePayload)
	return result
}
