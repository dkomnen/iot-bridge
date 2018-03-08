package thermometer

import (
	"bytes"
	"context"
	"encoding/binary"
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
	// if we did unit := 'c', the type of `unit` would default to rune, which is
	// 4 bytes long, and we don't want that, especially because the decoder will
	// expect 1 byte here.
	var unit byte = 'c'
	if v, ok := t.opts.Custom.Value(fahrenheit).(bool); ok && v {
		unit = 'f'
	}

	var low, high float64
	if v, ok := t.opts.Custom.Value(lowerBound).(float64); ok {
		low = v
	}
	if v, ok := t.opts.Custom.Value(higherBound).(float64); ok {
		high = v
	}

	return encode(
		t.opts.SerialNumber,
		unit,
		randomFloat64InRange(low, high),
		time.Now().Unix(),
	)
}

func encode(data ...interface{}) []byte {
	var encoded bytes.Buffer

	for _, v := range data {
		binary.Write(&encoded, binary.BigEndian, v)
	}

	return encoded.Bytes()
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
