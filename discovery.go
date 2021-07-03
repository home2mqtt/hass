package hass

import (
	"encoding/json"
	"errors"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type IConfig interface {
	GetComponent() string
	consume(consumer ConfigConsumer, nodeID string, objectID string)
}

func (*Cover) GetComponent() string {
	return "cover"
}

func (c *Cover) consume(consumer ConfigConsumer, nodeID string, objectID string) {
	consumer.ConsumeCover(c, nodeID, objectID)
}

func (*DInput) GetComponent() string {
	return "dinput"
}

func (c *DInput) consume(consumer ConfigConsumer, nodeID string, objectID string) {
	consumer.ConsumeDInput(c, nodeID, objectID)
}

func (*HVAC) GetComponent() string {
	return "hvac"
}

func (c *HVAC) consume(consumer ConfigConsumer, nodeID string, objectID string) {
	consumer.ConsumeHVAC(c, nodeID, objectID)
}

func (*Switch) GetComponent() string {
	return "switch"
}

func (c *Switch) consume(consumer ConfigConsumer, nodeID string, objectID string) {
	consumer.ConsumeSwitch(c, nodeID, objectID)
}

func (*Sensor) GetComponent() string {
	return "sensor"
}

func (c *Sensor) consume(consumer ConfigConsumer, nodeID string, objectID string) {
	consumer.ConsumeSensor(c, nodeID, objectID)
}

func (*Light) GetComponent() string {
	return "light"
}

func (c *Light) consume(consumer ConfigConsumer, nodeID string, objectID string) {
	consumer.ConsumeLight(c, nodeID, objectID)
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
}

func ProcessConfig(prefix string, topic string, payload []byte, consumer ConfigConsumer) (IConfig, error) {
	key := strings.TrimPrefix(topic, prefix+"/")
	keys := strings.Split(key, "/")
	if len(keys) != 4 {
		return nil, errors.New("Invalid topic: " + key)
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
	default:
		return nil, errors.New("Unknown component: " + topic)
	}
	err := json.Unmarshal(payload, config)
	if err != nil {
		return nil, err
	}
	config.consume(consumer, nodeid, objectid)
	return config, nil
}

func Discover(client mqtt.Client, consumer ConfigConsumer, errs chan error) {
	client.Subscribe(GetDiscoveryWildcard("homeassistant"), 0, func(client mqtt.Client, msg mqtt.Message) {
		_, err := ProcessConfig("homeassistant", msg.Topic(), msg.Payload(), consumer)
		if (err != nil) && (errs != nil) {
			errs <- err
		}
	}).Wait()
	client.Publish("discover", 0, false, "1").Wait()
}
