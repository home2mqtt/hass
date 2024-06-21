package hass

type ISensor[SensorEvent any] interface {
	Close() error
	Events() chan SensorEvent
}

type IField[Type any] interface {
	ISensor[Type]
	SetValue(value Type)
}

type IBoolField interface {
	IField[bool]
	Toggle()
}

type IEnumField interface {
	IField[string]
	List() []string
}

type ActionEvent struct {
	Action string
}

type IHVAC interface {
	Power() IField[bool]
	Mode() IEnumField
	Fan() IEnumField
	Swing() IEnumField
	TargetTemp() IField[float64]
	Temp() ISensor[float64]
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
	IField[int]
	Scale() int
}
