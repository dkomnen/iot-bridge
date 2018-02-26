package temp

import (
	"bytes"
	"context"
	"fmt"
	"time"

	"github.com/dkomnen/iot-bridge/cmd/valkyrie/device"
	"github.com/dkomnen/iot-bridge/mqtt"
)

const BrokerTopic = "TEMP"

type Temp struct {
	opts device.Options
	stop chan struct{}
}

func (t *Temp) Setup() error {
	if err := t.opts.Broker.Connect(); err != nil {
		return err
	}
	return nil
}

func (t *Temp) Run() error {
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

func (t *Temp) generateMessage() []byte {
	var buff bytes.Buffer
	unit := "celsius"
	if v, ok := t.opts.Custom.Value(fahrenheit).(bool); ok && v {
		unit = "fahrenheit"
	}

	var low, high float64
	if v, ok := t.opts.Custom.Value(lowerBound).(float64); ok {
		low = v
	}
	if v, ok := t.opts.Custom.Value(higherBound).(float64); ok {
		high = v
	}

	fmt.Fprintf(&buff, "%v:%s:%f", t.opts.SerialNumber, unit, randomFloat64InRange(low, high))

	return buff.Bytes()
}

func (t *Temp) Stop() error {
	t.stop <- struct{}{}
	close(t.stop)
	return nil
}

func (t *Temp) Options() device.Options {
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

	return &Temp{
		opts: defaults,
		stop: make(chan struct{}),
	}
}
