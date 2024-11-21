package bridge

import (
	"fmt"
	"log"

	"github.com/home2mqtt/hass"
)

func AttachSensor[T any](runtime hass.IPubSubRuntime, stateTopic string, valueTemplate string, sensor hass.ISensor[T]) {
	sensor.ReceiveEvent(func(value T) {
		runtime.Send(stateTopic, []byte(fmt.Sprint(value)))
	})
}

func AttachField[T any](runtime hass.IPubSubRuntime, stateTopic string, commandTopic string, field hass.IField[T], parseValue func(str string) (T, error)) {
	runtime.Receive(commandTopic, func(topic string, payload []byte) {
		value, err := parseValue(string(payload))
		if err != nil {
			log.Println(err)
			return
		}
		field.SetValue(value)
	})
	AttachSensor[T](runtime, stateTopic, "", field)
}
