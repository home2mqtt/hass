package proxy

import "github.com/home2mqtt/hass"

type button_impl struct {
	runtime hass.IPubSubRuntime
}

// ReceiveEvent implements hass.ISensor.
func (b *button_impl) ReceiveEvent(received func(event string)) {
	b.runtime.Receive("button", func(topic string, payload []byte) {
		received(string(payload))
	})
}

// SendEvent implements hass.ISensor.
func (b *button_impl) SendEvent(event string) {
	panic("unimplemented")
}

func NewButton(runtime hass.IPubSubRuntime, config *hass.Sensor) hass.ISensor[string] {
	return &button_impl{runtime: runtime}
}
