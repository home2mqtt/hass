package hass

type Device struct {
	Identifiers  []string `json:"identifiers,omitempty"`
	Connections  []string `json:"connections,omitempty"`
	Manufacturer string   `json:"manufacturer,omitempty"`
	Model        string   `json:"model,omitempty"`
	Name         string   `json:"name,omitempty"`
	SwVersion    string   `json:"sw_version,omitempty"`
}

type BasicConfig struct {
	Device   *Device `json:"device,omitempty"`
	UniqueID string  `json:"unique_id,omitempty"`
}

type Cover struct {
	BasicConfig
	CommandTopic   string `json:"command_topic,omitempty"`
	Name           string `json:"name,omitempty"`
	PositionTopic  string `json:"position_topic,omitempty"`
	PositionOpen   int    `json:"position_open"`
	PositionClosed int    `json:"position_closed"`
}

type DInput struct {
	BasicConfig
	Name       string `json:"name"`
	StateTopic string `json:"state_topic"`
}

// https://www.home-assistant.io/integrations/climate.mqtt/
type HVAC struct {
	BasicConfig
	Name string `json:"name,omitempty"`

	// A template to render the value received on the action_topic with.
	ActionTemplate string `json:"action_template,omitempty"`
	// The MQTT topic to subscribe for changes of the current action. If this is set, the climate graph uses the value received as data source. Valid values: off, heating, cooling, drying, idle, fan.
	ActionTopic string `json:"action_topic,omitempty"`

	// A template with which the value received on current_humidity_topic will be rendered.
	CurrentHumidityTemplate string `json:"current_humidity_template,omitempty"`
	// The MQTT topic on which to listen for the current humidity. A "None" value received will reset the current humidity. Empty values (''') will be ignored.
	CurrentHumidityTopic string `json:"current_humidity_topic,omitempty"`

	// A template with which the value received on current_temperature_topic will be rendered.
	CurrentTemperatureTemplate string `json:"current_temperature_template,omitempty"`
	// The MQTT topic on which to listen for the current temperature. A "None" value received will reset the current temperature. Empty values (''') will be ignored.
	CurrentTemperatureTopic string `json:"current_temperature_topic,omitempty"`

	// A template to render the value sent to the fan_mode_command_topic with.
	FanModeCommandTemplate string `json:"fan_mode_command_template,omitempty"`
	// The MQTT topic to publish commands to change the fan mode.
	FanModeCommandTopic string `json:"fan_mode_command_topic,omitempty"`
	// A template to render the value received on the fan_mode_state_topic with.
	FanModeStateTemplate string `json:"fan_mode_state_template,omitempty"`
	// The MQTT topic to subscribe for changes of the HVAC fan mode. If this is not set, the fan mode works in optimistic mode (see below).
	FanModeStateTopic string `json:"fan_mode_state_topic,omitempty"`
	//A list of supported fan modes. Default: [“auto”, “low”, “medium”, “high”]
	FanModes []string `json:"fan_modes,omitempty"`

	// The minimum target humidity percentage that can be set.
	MaxHumidity string `json:"max_humidity,omitempty"`
	// Maximum set point available. The default value depends on the temperature unit, and will be 35°C or 95°F.
	MaxTemp float64 `json:"max_temp,omitempty"`
	// The maximum target humidity percentage that can be set.
	MinHumidity string `json:"min_humidity,omitempty"`
	// Minimum set point available. The default value depends on the temperature unit, and will be 7°C or 44.6°F.
	MinTemp float64 `json:"min_temp,omitempty"`

	// A template to render the value sent to the mode_command_topic with.
	ModeCommandTemplate string `json:"mode_command_template,omitempty"`
	// The MQTT topic to publish commands to change the HVAC operation mode.
	ModeCommandTopic string `json:"mode_command_topic,omitempty"`
	// A template to render the value received on the mode_state_topic with.
	ModeStateTemplate string `json:"mode_state_template,omitempty"`
	// The MQTT topic to subscribe for changes of the HVAC operation mode. If this is not set, the operation mode works in optimistic mode (see below).
	ModeStateTopic string `json:"mode_state_topic,omitempty"`
	// A list of supported modes. Needs to be a subset of the default values. Default: [“auto”, “off”, “cool”, “heat”, “dry”, “fan_only”]
	Modes []string `json:"modes,omitempty"`

	// Flag that defines if the climate works in optimistic mode Default: true if no state topic defined, else false.
	Optmistic bool `json:"optimistic,omitempty"`

	// The payload sent to turn off the device. (optional, default: OFF)
	PayloadOff any `json:"payload_off,omitempty"`
	// The payload sent to turn the device on. (optional, default: ON)
	PayloadOn any `json:"payload_on,omitempty"`

	// A template to render the value sent to the power_command_topic with. The value parameter is the payload set for payload_on or payload_off.
	PowerCommandTemplate string `json:"power_command_template,omitempty"`
	// The MQTT topic to publish commands to change the HVAC power state. Sends the payload configured with payload_on if the climate is turned on via the climate.turn_on, or the payload configured with payload_off if the climate is turned off via the climate.turn_off service. Note that optimistic mode is not supported through climate.turn_on and climate.turn_off services. When called, these services will send a power command to the device but will not optimistically update the state of the climate entity. The climate device should report its state back via mode_state_topic.
	PowerCommandTopic string `json:"power_command_topic,omitempty"`

	// The desired precision for this device. Can be used to match your actual thermostat’s precision. Supported values are 0.1, 0.5 and 1.0. Default: 0.1 for Celsius and 1.0 for Fahrenheit.
	Precision float64 `json:"precision,omitempty"`

	// Defines a template to generate the payload to send to preset_mode_command_topic.
	PresetModeCommandTemplate string `json:"preset_mode_command_template,omitempty"`
	// The MQTT topic to publish commands to change the preset mode.
	PresetModeCommandTopic string `json:"preset_mode_command_topic,omitempty"`
	// The MQTT topic subscribed to receive climate speed based on presets. When preset ‘none’ is received or None the preset_mode will be reset.
	PresetModeStateTopic string `json:"preset_mode_state_topic,omitempty"`
	// List of preset modes this climate is supporting. Common examples include eco, away, boost, comfort, home, sleep and activity.
	PresetModes []string `json:"preset_modes,omitempty"`

	// The MQTT topic to publish commands to change the swing mode.
	SwingModeCommandTopic string `json:"swing_mode_command_topic,omitempty"`
	// The MQTT topic to subscribe for changes of the HVAC swing mode. If this is not set, the swing mode works in optimistic mode (see below).
	SwingModeStateTopics string `json:"swing_mode_state_topic,omitempty"`
	// A list of supported swing modes. default: [“on”, “off”]
	SwingModes []string `json:"swing_modes,omitempty"`

	// The MQTT topic to publish commands to change the target humidity.
	TargetHumidityCommandTopic string `json:"target_humidity_command_topic,omitempty"`
	// The MQTT topic subscribed to receive the target humidity. If this is not set, the target humidity works in optimistic mode (see below). A "None" value received will reset the target humidity. Empty values (''') will be ignored.
	TargetHumidityStateTopic string `json:"target_humidity_state_topic,omitempty"`

	// The MQTT topic to publish commands to change the target temperature.
	TemperatureCommandTopic string `json:"temperature_command_topic,omitempty"`
	// The MQTT topic to subscribe for changes in the target temperature. If this is not set, the target temperature works in optimistic mode (see below). A "None" value received will reset the temperature set point. Empty values (''') will be ignored.
	TemperatureStateTopic string `json:"temperature_state_topic,omitempty"`

	TempStep float64 `json:"temp_step,omitempty"`
}

type Light struct {
	BasicConfig
	CommandTopic           string   `json:"command_topic,omitempty"`
	Name                   string   `json:"name,omitempty"`
	Brightness             bool     `json:"brightness,omitempty"`
	BrightnessCommandTopic string   `json:"brightness_command_topic,omitempty"`
	BrightnessScale        int32    `json:"brightness_scale"`
	BrightnessStateTopic   string   `json:"brightness_state_topic,omitempty"`
	ColorMode              bool     `json:"color_mode,omitempty"`
	SupportedColorModes    []string `json:"supported_color_modes,omitempty"`
	OnCommandType          string   `json:"on_command_type,omitempty"`
	StateTopic             string   `json:"state_topic,omitempty"`
	Effect                 bool     `json:"effect,omitempty"`
	EffectList             []string `json:"effect_list,omitempty"`
	HS                     bool     `json:"hs,omitempty"`
	JsonAttributesTopic    string   `json:"json_attributes_topic,omitempty"`
	Schema                 string   `json:"schema,omitempty"`
	XY                     bool     `json:"xy,omitempty"`
}

// https://www.home-assistant.io/integrations/sensor.mqtt/
type Sensor struct {
	BasicConfig
	Name              string `json:"name,omitempty"`
	UnitOfMeasurement string `json:"unit_of_measurement"`
	Topic             string `json:"state_topic"`
	Icon              string `json:"icon,omitempty"`
	StateClass        string `json:"state_class,omitempty"`
	ValueTemplate     string `json:"value_template,omitempty"`
	DeviceClass       string `json:"device_class,omitempty"`
}

// https://www.home-assistant.io/integrations/switch.mqtt/
type Switch struct {
	BasicConfig
	CommandTopic  string `json:"command_topic,omitempty"`
	Name          string `json:"name,omitempty"`
	StateTopic    string `json:"state_topic,omitempty"`
	ValueTemplate string `json:"value_template,omitempty"`
	PayLoadOff    any    `json:"payload_off,omitempty"`
	PayLoadOn     any    `json:"payload_on,omitempty"`
}

type BinarySensor struct {
	BasicConfig
	Name          string `json:"name,omitempty"`
	ValueTemplate string `json:"value_template,omitempty"`
	DeviceClass   string `json:"device_class,omitempty"`
	Topic         string `json:"state_topic"`
	Icon          string `json:"icon,omitempty"`
	PayLoadOff    any    `json:"payload_off,omitempty"`
	PayLoadOn     any    `json:"payload_on,omitempty"`
}
