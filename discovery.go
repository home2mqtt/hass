package hass

import (
	"encoding/json"
	"errors"
	"strings"
)

type IConfig interface {
	GetComponent() string
	Consume(consumer ConfigConsumer, nodeID string, objectID string)
	GetBasic() *BasicConfig
}

type DiscoveredConfig struct {
	Topic   string
	Payload []byte
}

func (c *BasicConfig) GetBasic() *BasicConfig {
	return c
}

func (*Cover) GetComponent() string {
	return "cover"
}

func (c *Cover) Consume(consumer ConfigConsumer, nodeID string, objectID string) {
	consumer.ConsumeCover(c, nodeID, objectID)
}

func (*DInput) GetComponent() string {
	return "dinput"
}

func (c *DInput) Consume(consumer ConfigConsumer, nodeID string, objectID string) {
	consumer.ConsumeDInput(c, nodeID, objectID)
}

func (*HVAC) GetComponent() string {
	return "hvac"
}

func (c *HVAC) Consume(consumer ConfigConsumer, nodeID string, objectID string) {
	consumer.ConsumeHVAC(c, nodeID, objectID)
}

func (*Switch) GetComponent() string {
	return "switch"
}

func (c *Switch) Consume(consumer ConfigConsumer, nodeID string, objectID string) {
	consumer.ConsumeSwitch(c, nodeID, objectID)
}

func (*Sensor) GetComponent() string {
	return "sensor"
}

func (c *Sensor) Consume(consumer ConfigConsumer, nodeID string, objectID string) {
	consumer.ConsumeSensor(c, nodeID, objectID)
}

func (*Light) GetComponent() string {
	return "light"
}

func (c *Light) Consume(consumer ConfigConsumer, nodeID string, objectID string) {
	consumer.ConsumeLight(c, nodeID, objectID)
}

func (*BinarySensor) GetComponent() string {
	return "binary_sensor"
}

func (c *BinarySensor) Consume(consumer ConfigConsumer, nodeID string, objectID string) {
	consumer.ConsumeBinarySensor(c, nodeID, objectID)
}

func GetDiscoveryTopic(c IConfig, prefix string, nodeID string, objectID string) string {
	return prefix + "/" + c.GetComponent() + "/" + nodeID + "/" + objectID + "/config"
}

func GetPayload(c IConfig) ([]byte, error) {
	return json.Marshal(c)
}

func GetDiscoveryWildcard(prefix string) string {
	return prefix + "/+/+/+/config"
}

type ConfigConsumer interface {
	ConsumeCover(c *Cover, nodeID string, objectID string)
	ConsumeDInput(c *DInput, nodeID string, objectID string)
	ConsumeHVAC(c *HVAC, nodeID string, objectID string)
	ConsumeSwitch(c *Switch, nodeID string, objectID string)
	ConsumeSensor(c *Sensor, nodeID string, objectID string)
	ConsumeLight(c *Light, nodeID string, objectID string)
	ConsumeBinarySensor(c *BinarySensor, nodeID string, objectID string)
}

func ProcessConfig(prefix string, confData DiscoveredConfig) (IConfig, string, string, error) {
	key := strings.TrimPrefix(confData.Topic, prefix+"/")
	keys := strings.Split(key, "/")
	if len(keys) != 4 {
		return nil, "", "", errors.New("Invalid topic: " + key)
	}
	component := keys[0]
	nodeid := keys[1]
	objectid := keys[2]
	var config IConfig
	switch component {
	case "cover":
		config = &Cover{}
	case "dinput":
		config = &DInput{}
	case "hvac":
		config = &HVAC{}
	case "switch":
		config = &Switch{}
	case "sensor":
		config = &Sensor{}
	case "light":
		config = &Light{}
	case "binary_sensor":
		config = &BinarySensor{}
	default:
		return nil, "", "", errors.New("Unknown component: " + confData.Topic)
	}
	err := json.Unmarshal(confData.Payload, config)
	if err != nil {
		return nil, "", "", err
	}
	return config, nodeid, objectid, nil
}

func Discover(client IPubSubRuntime, prefix string) (chan DiscoveredConfig, error) {
	result := make(chan DiscoveredConfig)

	client.Receive(GetDiscoveryWildcard(prefix), func(topic string, payload []byte) {
		result <- DiscoveredConfig{
			Topic:   topic,
			Payload: payload,
		}
	})

	client.Send("discover", []byte("1"))
	return result, nil
}

func ConsumeDiscoveredConfigs(client IPubSubRuntime, consumer ConfigConsumer, errs chan error) {
	client.Receive(GetDiscoveryWildcard("homeassistant"), func(topic string, payload []byte) {
		config, nodeid, objectid, err := ProcessConfig("homeassistant", DiscoveredConfig{
			Topic:   topic,
			Payload: payload,
		})
		if (err != nil) && (errs != nil) {
			errs <- err
		} else {
			config.Consume(consumer, nodeid, objectid)
		}
	})
	client.Send("discover", []byte("1"))
}
