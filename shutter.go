package hass

type shutter struct {
	baseActuator
	basic bool
}

func BasicShutter(runtime IPubSubRuntime, conf *Cover) IShutter {
	result := &shutter{
		basic: true,
	}
	result.init(runtime, conf.CommandTopic)
	return result
}

func Shutter(runtime IPubSubRuntime, conf *Cover) IShutter {
	result := &shutter{}
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
