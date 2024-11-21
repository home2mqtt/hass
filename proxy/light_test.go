package proxy_test

import (
	"testing"

	"github.com/home2mqtt/hass"
	"github.com/home2mqtt/hass/mock"
	"github.com/home2mqtt/hass/proxy"
)

func TestLightIsOn(t *testing.T) {
	config := &hass.Light{
		StateTopic: "test/state",
	}
	mockPubSub := mock.NewMockPubSub()
	var getted bool
	mockPubSub.Receive(config.StateTopic+"/get", func(topic string, payload []byte) {
		go mockPubSub.Send(config.StateTopic, []byte(`{"state":"on"}`))
		getted = true
	})

	light := proxy.NewLight(mockPubSub, config)
	events := make(chan bool)

	statesensor := light.State()
	statesensor.ReceiveEvent(func(event bool) {
		events <- event
	})
	if !getted {
		t.Fail()
	}
	value := <-events
	if !value {
		t.Fail()
	}
}
