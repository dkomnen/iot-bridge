package broker

type Broker interface {
	Connect() error
	Disconnect() error
	Publish(topic string, msg []byte, opts ...PublishOption) error
	Subscribe(topic string, handler Handler, opts ...SubscribeOption) error
	Options() Options
}

type Handler func([]byte) error
