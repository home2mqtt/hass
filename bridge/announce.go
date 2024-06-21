package bridge

import (
	"encoding/json"

	"github.com/home2mqtt/hass"
)

// <discovery_prefix>/<component>/[<node_id>/]<object_id>/config
func AnnounceDevice(client hass.IPubSubRuntime, prefix string, nodeid string, objectid string, device hass.IConfig) error {
	c, err := json.Marshal(device)
	if err != nil {
		return err
	}
	topic := prefix + "/" + device.GetComponent() + "/" + nodeid + "/" + objectid + "/config"
	client.SendRetained(topic, c)
	return nil
}
