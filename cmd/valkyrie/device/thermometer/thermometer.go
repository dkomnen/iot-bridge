package thermometer

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/dkomnen/iot-bridge/broker/mqtt"
	"github.com/dkomnen/iot-bridge/cmd/valkyrie/device"
)

const BrokerTopic = "THERMOMETER"

type Thermometer struct {
	opts      device.Options
	isRunning bool
	stop      chan struct{}
}

func (t *Thermometer) Setup() error {
	if err := t.opts.Broker.Connect(); err != nil {
		return err
	}
	return nil
}

func (t *Thermometer) Run() error {
	t.isRunning = true
	defer func() { t.isRunning = false }()

	tick := time.NewTicker(t.opts.Interval)
	for {
		select {
		case <-tick.C:
			if err := t.opts.Broker.Publish(
				BrokerTopic,
				t.generateMessage(),
			); err != nil {
				return err
			}
		case <-t.stop:
			return nil
		}
	}
	return nil
}

func (t *Thermometer) generateMessage() []byte {
	var buff bytes.Buffer
	unit := "c"
	if v, ok := t.opts.Custom.Value(fahrenheit).(bool); ok && v {
		unit = "f"
	}

	var low, high float64
	if v, ok := t.opts.Custom.Value(lowerBound).(float64); ok {
		low = v
	}
	if v, ok := t.opts.Custom.Value(higherBound).(float64); ok {
		high = v
	}

	fmt.Fprintf(&buff, "%s%s%5.3f", t.opts.SerialNumber, unit, randomFloat64InRange(low, high))

	return buff.Bytes()
}

func (t *Thermometer) Stop() error {
	if t.isRunning {
		t.stop <- struct{}{}
	}
	close(t.stop)
	return nil
}

func (t *Thermometer) Options() device.Options {
	return t.opts
}

func New(opts ...device.Option) device.Device {
	defaults := device.Options{
		SerialNumber: [32]byte{0},
		Broker:       mqtt.New(),
		Interval:     time.Second,
		Custom:       context.Background(),
	}

	for _, o := range opts {
		o(&defaults)
	}

	return &Thermometer{
		opts: defaults,
		stop: make(chan struct{}),
	}
}
