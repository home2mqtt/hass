package proxy

import (
	"fmt"

	"github.com/home2mqtt/hass"
)

type contact_impl struct {
	runtime  hass.IPubSubRuntime
	onValue  string
	offValue string
}

// ReceiveEvent implements hass.ISensor.
func (c *contact_impl) ReceiveEvent(received func(event bool)) {
	c.runtime.Receive("contact", func(topic string, payload []byte) {
		value := string(payload)
		on := c.onValue == value
		received(on)
	})
}

// SendEvent implements hass.ISensor.
func (c *contact_impl) SendEvent(event bool) {
	var value string
	if event {
		value = c.onValue
	} else {
		value = c.offValue
	}
	c.runtime.Send("contact", []byte(value))
}

func NewContact(runtime hass.IPubSubRuntime, config *hass.BinarySensor) hass.ISensor[bool] {
	var onValue string
	if config.PayLoadOn == nil {
		onValue = "ON"
	} else {
		onValue = fmt.Sprintf("%v", config.PayLoadOn)
	}
	var offValue string
	if config.PayLoadOff == nil {
		offValue = "OFF"
	} else {
		offValue = fmt.Sprintf("%v", config.PayLoadOff)
	}
	return &contact_impl{runtime: runtime, onValue: onValue, offValue: offValue}
}
