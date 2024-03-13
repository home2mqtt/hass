package hass

import (
	"encoding/json"
	"log"

	"github.com/noirbizarre/gonja"
)

func GenericSensor(context IPubSubRuntime, sensor *Sensor) ISensor[string] {
	values := ParseSensorValue(context, sensor.Topic, sensor.ValueTemplate, func(s string) (string, error) {
		return s, nil
	})
	result := &BaseSensor[string]{
		events: values,
	}

	return result
}

func ParseSensorValue[T any](context IPubSubRuntime, Topic string, ValueTemplate string, parse func(string) (T, error)) chan T {
	result := make(chan T)

	tpl, err := gonja.FromString(ValueTemplate)
	if err != nil {
		log.Println(err)
	}
	context.Receive(Topic, func(topic string, payload []byte) {
		var data map[string]interface{}
		err := json.Unmarshal(payload, &data)
		if err != nil {
			log.Println(err)
			return
		}
		context := map[string]interface{}{"value_json": data}
		value, err := tpl.Execute(context)
		if err != nil {
			log.Println(err)
			return
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
	return result
}
