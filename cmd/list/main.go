package main

import (
	"log"
	"time"

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

func main() {
	opts := mqtt.NewClientOptions().AddBroker("tcp://192.168.0.1:1883").SetAutoReconnect(true)
	client := mqtt.NewClient(opts)
	token := client.Connect()
	token.Wait()
	err := token.Error()
	if err != nil {
		log.Println("Connection failed: ", token.Error())
	}

	var consumer hass.ConfigConsumer = &logConsumer{}

	client.Subscribe(hass.GetDiscoveryWildcard("homeassistant"), 0, func(client mqtt.Client, msg mqtt.Message) {
		_, err := hass.ProcessConfig("homeassistant", msg.Topic(), msg.Payload(), consumer)
		if err != nil {
			log.Println(err)
		}
	})
	client.Publish("discover", 0, false, 0)

	time.Sleep(10 * time.Second)
	client.Disconnect(100)
}
