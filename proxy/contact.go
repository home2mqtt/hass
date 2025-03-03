package proxy

import (
	"fmt"

	"github.com/home2mqtt/hass"
)

type contact_impl struct {
	runtime hass.IPubSubRuntime
	onValue string
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
	panic("unimplemented")
}

func NewContact(runtime hass.IPubSubRuntime, config *hass.BinarySensor) hass.ISensor[bool] {
	var onValue string
	if config.PayLoadOn == nil {
		onValue = "ON"
	} else {
		onValue = fmt.Sprintf("%v", config.PayLoadOn)
	}
	return &contact_impl{runtime: runtime, onValue: onValue}
}
