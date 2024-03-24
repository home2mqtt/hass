package hass

type IPubSubRuntime interface {
	Send(topic string, payload []byte) error
	Receive(topicpattern string, callback func(topic string, payload []byte))
}

type IConnecter interface {
	Connect() error
	Disconnect()
	ConnectionState() chan bool
}

type IPubSubConnecterRuntime interface {
	IPubSubRuntime
	IConnecter
}
