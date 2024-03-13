package hass

type ISensor[SensorEvent any] interface {
	Close() error
	Events() chan SensorEvent
}

type ActionEvent struct {
	Action string
}

type IHVAC[StateType any] interface {
	Stop() (StateType, error)
	Restart(StateType) error
	State() StateType
	IsOn(StateType) bool
}

type IBasicShutter interface {
	Open()
	Close()
	Stop()
}

type IShutter interface {
	IBasicShutter
	OpenOrStop()
	CloseOrStop()
}

type ISwitch interface {
	State() ISensor[bool]
	Toggle()
	Set(on bool)
}

type ILight interface {
	ISwitch
	Brightness() IIntSettable
	ColorTemp() IIntSettable
}

type IIntSettable interface {
	SetValue(value int)
	Scale() int
	Values() chan int
}
