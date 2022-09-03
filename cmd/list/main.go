package main

import (
	"log"

	"github.com/balazsgrill/hass"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type logConsumer struct{}

func (*logConsumer) ConsumeCover(c *hass.Cover, nodeID string, objectID string) {
	log.Printf("Cover: %s, %s", nodeID, objectID)
}
func (*logConsumer) ConsumeDInput(c *hass.DInput, nodeID string, objectID string) {
	log.Printf("DInput: %s, %s", nodeID, objectID)
}
func (*logConsumer) ConsumeHVAC(c *hass.HVAC, nodeID string, objectID string) {
	log.Printf("HVAC: %s, %s", nodeID, objectID)
}
func (*logConsumer) ConsumeSwitch(c *hass.Switch, nodeID string, objectID string) {
	log.Printf("Switch: %s, %s", nodeID, objectID)
}
func (*logConsumer) ConsumeSensor(c *hass.Sensor, nodeID string, objectID string) {
	log.Printf("Sensor: %s, %s", nodeID, objectID)
}
func (*logConsumer) ConsumeLight(c *hass.Light, nodeID string, objectID string) {
	log.Printf("Light: %s, %s", nodeID, objectID)
}
func (*logConsumer) ConsumeBinarySensor(c *hass.BinarySensor, nodeID string, objectID string) {
	log.Printf("Binary Sensor: %s, %s", nodeID, objectID)
}

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://192.168.0.1:1883").SetAutoReconnect(true)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()
	err := token.Error()
	if err != nil {
		log.Println("Connection failed: ", token.Error())
	}
	defer client.Disconnect(100)

	var consumer hass.ConfigConsumer = &logConsumer{}

	errs := make(chan error)
	hass.ConsumeDiscoveredConfigs(client, consumer, errs)

	for err = range errs {
		log.Println(err)
	}
}
