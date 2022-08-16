package hass

import (
	"encoding/json"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

// <discovery_prefix>/<component>/[<node_id>/]<object_id>/config
func AnnounceDevice(client mqtt.Client, prefix string, nodeid string, objectid string, device IConfig) error {
	c, err := json.Marshal(device)
	if err != nil {
		return err
	}
	topic := prefix + "/" + device.GetComponent() + "/" + nodeid + "/" + objectid + "/config"
	client.Publish(topic, 0, true, c)
	return nil
}
