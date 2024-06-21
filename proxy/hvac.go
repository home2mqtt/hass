package proxy

import "github.com/home2mqtt/hass"

type hvac struct {
	client     hass.IPubSubRuntime
	power      hass.IBoolField
	mode       hass.IEnumField
	fan        hass.IEnumField
	swing      hass.IEnumField
	targetTemp hass.IField[float64]
	temp       hass.ISensor[float64]
}

// NewHVAC returns a new HVAC.
func NewHVAC(client hass.IPubSubRuntime, config *hass.HVAC) hass.IHVAC {
	device := &hvac{
		client: client,
	}
	device.power = NewBoolField(client, "", config.PowerCommandTopic)
	device.mode = NewEnumField(client, config.ModeStateTopic, config.ModeCommandTopic, config.Modes)
	device.fan = NewEnumField(client, config.FanModeStateTopic, config.FanModeCommandTopic, config.FanModes)
	device.swing = NewEnumField(client, config.SwingModeStateTopics, config.SwingModeCommandTopic, config.SwingModes)
	device.targetTemp = NewFloatField(client, "", config.TemperatureStateTopic, config.TemperatureCommandTopic)
	device.temp = NewFloatSensor(client, config.CurrentTemperatureTemplate, config.CurrentTemperatureTopic)
	return device
}

// Fan implements IHVAC.
func (h *hvac) Fan() hass.IEnumField {
	return h.fan
}

// Mode implements IHVAC.
func (h *hvac) Mode() hass.IEnumField {
	return h.mode
}

// Power implements IHVAC.
func (h *hvac) Power() hass.IField[bool] {
	return h.power
}

// Swing implements IHVAC.
func (h *hvac) Swing() hass.IEnumField {
	return h.swing
}

// TargetTemp implements IHVAC.
func (h *hvac) TargetTemp() hass.IField[float64] {
	return h.targetTemp
}

// Temp implements IHVAC.
func (h *hvac) Temp() hass.ISensor[float64] {
	return h.temp
}

var _ hass.IHVAC = &hvac{}
