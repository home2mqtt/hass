package bridge

import (
	"log"
	"strings"

	"github.com/home2mqtt/hass"
)

type PropertyContext struct {
	hass.IPubSubRuntime
	Base string
	Id   string
}

type IProperty[T any] interface {
	StateTopic() string
	CommandTopic() string
	NotifyState(value T)
	OnCommand(callback func(value T))
}

func (pc *PropertyContext) DefineString(name string) IProperty[string] {
	return &stringProperty{
		property: property{
			PropertyContext: pc,
			name:            name,
		},
	}
}

func (pc *PropertyContext) DefineFloat(name string) IProperty[float64] {
	return &floatProperty{
		property: property{
			PropertyContext: pc,
			name:            name,
		},
	}
}

func (pc *PropertyContext) DefineOnOff(name string) IProperty[bool] {
	return &onOffProperty{
		property: property{
			PropertyContext: pc,
			name:            name,
		},
		onValue:  "ON",
		offValue: "OFF",
	}
}

type property struct {
	*PropertyContext
	name string
}

func (p *property) StateTopic() string {
	return strings.Join([]string{p.Base, p.Id, p.name}, "/")
}

func (p *property) CommandTopic() string {
	return strings.Join([]string{p.Base, p.Id, p.name, "set"}, "/")
}

type stringProperty struct {
	property
}

func (p *stringProperty) NotifyState(value string) {
	hass.SendString(p, p.StateTopic(), value)
}

func (p *stringProperty) OnCommand(callback func(value string)) {
	hass.ReceiveString(p, p.CommandTopic(), func(topic, payload string) {
		callback(payload)
	})
}

type floatProperty struct {
	property
}

func (p *floatProperty) NotifyState(value float64) {
	hass.SendFloat(p, p.StateTopic(), value)
}

func (p *floatProperty) OnCommand(callback func(value float64)) {
	hass.ReceiveFloat(p, p.CommandTopic(), func(topic string, payload float64, err error) {
		if err == nil {
			callback(payload)
		} else {
			log.Printf("Float value error received on %s: %v\n", topic, err)
		}
	})
}

type onOffProperty struct {
	property
	onValue  string
	offValue string
}

func (p *onOffProperty) NotifyState(value bool) {
	if value {
		hass.SendString(p, p.StateTopic(), p.onValue)
	} else {
		hass.SendString(p, p.StateTopic(), p.offValue)
	}
}

func (p *onOffProperty) OnCommand(callback func(value bool)) {
	hass.ReceiveString(p, p.CommandTopic(), func(topic, payload string) {
		if strings.EqualFold(payload, p.onValue) {
			callback(true)
		} else if strings.EqualFold(payload, p.offValue) {
			callback(false)
		} else {
			log.Printf("Invalid value received on %s: %s\n", topic, payload)
		}
	})
}
