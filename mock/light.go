package mock

import "github.com/home2mqtt/hass"

type MockLight struct {
	*MockSwitch
}

var _ hass.ILight = (*MockLight)(nil)

func (l *MockLight) Brightness() hass.IIntSettable {
	return nil
}

func (l *MockLight) ColorTemp() hass.IIntSettable {
	return nil
}

func NewMockLight() *MockLight {
	return &MockLight{
		MockSwitch: NewMockSwitch(),
	}
}
