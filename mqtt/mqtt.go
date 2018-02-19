package mqtt

import (
	"context"
	"sync"

	"github.com/dkomnen/iot-bridge/broker"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type mqttBroker struct {
	sync.Mutex
	client  mqtt.Client
	options broker.Options
}

func (m *mqttBroker) Connect() error {
	m.Lock()
	defer m.Unlock()

	if t := m.client.Connect(); t.Wait() && t.Error() != nil {
		return t.Error()
	}

	return nil
}

func (m *mqttBroker) Disconnect() error {
	m.Lock()
	defer m.Unlock()

	m.client.Disconnect(0)

	return nil
}

func (m *mqttBroker) Publish(topic string, msg []byte, opts ...broker.PublishOption) error {
	m.Lock()
	defer m.Unlock()

	options := broker.PublishOptions{Custom: context.Background()}
	for _, o := range opts {
		o(&options)
	}

	var qos byte
	if v, ok := options.Custom.Value(publishQoS).(byte); ok {
		qos = v
	}

	var retained bool
	if v, ok := options.Custom.Value(publishRetained).(bool); ok {
		retained = v
	}

	if t := m.client.Publish(topic, qos, retained, msg); t.Wait() && t.Error() != nil {
		return t.Error()
	}

	return nil
}

func (m *mqttBroker) Subscribe(topic string, handler broker.Handler, opts ...broker.SubscribeOption) error {
	m.Lock()
	defer m.Unlock()

	options := broker.SubscribeOptions{Custom: context.Background()}
	for _, o := range opts {
		o(&options)
	}

	var qos byte
	if v, ok := options.Custom.Value(subscribeQoS).(byte); ok {
		qos = v
	}

	callback := func(c mqtt.Client, msg mqtt.Message) {
		handler(msg.Payload())
	}

	if t := m.client.Subscribe(topic, qos, callback); t.Wait() && t.Error() != nil {
		return t.Error()
	}

	return nil
}

func (m *mqttBroker) Options() broker.Options {
	m.Lock()
	defer m.Unlock()

	return m.options
}

func getMQTTOpts(opts broker.Options) *mqtt.ClientOptions {
	mqttClientOpts := mqtt.NewClientOptions()
	mqttClientOpts.AddBroker(opts.Address)
	if opts.TLSConfig != nil {
		mqttClientOpts.SetTLSConfig(opts.TLSConfig)
	}

	setAllCustomClientOptions(mqttClientOpts, opts.Custom)

	return mqttClientOpts
}

func New(options ...broker.Option) broker.Broker {
	defaultOpts := broker.Options{
		Address:   "tcp://localhost:1883",
		TLSConfig: nil,
		Custom:    context.Background(),
	}

	for _, o := range options {
		o(&defaultOpts)
	}

	mqttOpts := getMQTTOpts(defaultOpts)

	return &mqttBroker{
		client:  mqtt.NewClient(mqttOpts),
		options: defaultOpts,
	}
}
