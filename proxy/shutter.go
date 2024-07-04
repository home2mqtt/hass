package proxy

import (
	"strconv"

	"github.com/home2mqtt/hass"
)

type shutter struct {
	baseActuator
	basic         bool
	positionrange int
	position      hass.ISensor[int]
}

// Position implements hass.IShutter.
func (s *shutter) Position() hass.ISensor[int] {
	return s.position
}

// Range implements hass.IShutter.
func (s *shutter) Range() int {
	return s.positionrange
}

func BasicShutter(runtime hass.IPubSubRuntime, conf *hass.Cover) hass.IShutter {
	result := &shutter{
		basic: true,
	}
	result.init(runtime, conf.CommandTopic)
	return result
}

func Shutter(runtime hass.IPubSubRuntime, conf *hass.Cover) hass.IShutter {
	result := &shutter{}
	if conf.PositionOpen != 0 {
		result.positionrange = conf.PositionOpen - conf.PositionClosed
	} else {
		result.positionrange = 100
	}
	if conf.PositionTopic != "" {
		result.position = hass.ChanToSensor(ParseSensorValue(runtime, conf.PositionTopic, conf.PositionTemplate, func(s string) (int, error) {
			return strconv.Atoi(s)
		}))
	}
	result.init(runtime, conf.CommandTopic)
	return result
}

func (s *shutter) Open() {
	s.send("OPEN")
}

func (s *shutter) Close() {
	s.send("CLOSE")
}

func (s *shutter) Stop() {
	s.send("STOP")
}

func (s *shutter) OpenOrStop() {
	if s.basic {
		s.Open()
	} else {
		s.send("OPENORSTOP")
	}
}

func (s *shutter) CloseOrStop() {
	if s.basic {
		s.Close()
	} else {
		s.send("CLOSEORSTOP")
	}
}
