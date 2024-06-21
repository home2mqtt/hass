package bridge

import (
	"fmt"

	"github.com/home2mqtt/hass"
)

func AttachSensor[T any](runtime hass.IPubSubRuntime, stateTopic string, valueTemplate string, sensor hass.ISensor[T]) {
	go func() {
		for v := range sensor.Events() {
			runtime.Send(stateTopic, []byte(fmt.Sprint(v)))
		}
	}()
}

func AttachField[T any](runtime hass.IPubSubRuntime, stateTopic string, commandTopic string, field hass.IField[T]) {
	runtime.Receive(commandTopic, func(topic string, payload []byte) {
		field.SetValue(field.ParseValue(payload))
	})
	AttachSensor[T](runtime, stateTopic, "", field)
}
