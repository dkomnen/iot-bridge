package broker

type Broker interface {
	Publish(topic string, msg []byte, opts ...PublishOption) error
	Subscribe(topic string, handler Handler, opts ...SubscribeOption) error
	Connect() error
	Disconnect() error
	Options() Options
}

type Handler func([]byte) error

var (
	defaultBroker = newMQTTBroker()
)

func NewBroker(options ...Option) Broker {
	return newMQTTBroker(options)
}
