package hass

import (
	"encoding/json"
)

// <discovery_prefix>/<component>/[<node_id>/]<object_id>/config
func AnnounceDevice(client IPubSubRuntime, prefix string, nodeid string, objectid string, device IConfig) error {
	c, err := json.Marshal(device)
	if err != nil {
		return err
	}
	topic := prefix + "/" + device.GetComponent() + "/" + nodeid + "/" + objectid + "/config"
	client.Send(topic, c)
	return nil
}
