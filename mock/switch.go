package mock

import "github.com/home2mqtt/hass"

type MockSwitch struct {
	hass.ISensor[bool]
	StateValue bool
}

func NewMockSwitch() *MockSwitch {
	return &MockSwitch{
		ISensor: hass.NewSensor[bool](),
	}
}

var _ hass.ISwitch = &MockSwitch{}

func (ms *MockSwitch) State() hass.ISensor[bool] {
	return ms
}
func (ms *MockSwitch) Set(value bool) {
	ms.StateValue = value
	ms.ISensor.SendEvent(ms.StateValue)
}

func (ms *MockSwitch) Toggle() {
	ms.StateValue = !ms.StateValue
	ms.ISensor.SendEvent(ms.StateValue)
}
