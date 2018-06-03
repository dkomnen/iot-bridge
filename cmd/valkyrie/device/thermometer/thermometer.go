package thermometer

import (
	"bytes"
	"context"
	"encoding/binary"
	"time"
	"fmt"

	"github.com/dkomnen/iot-bridge/broker/mqtt"
	"github.com/dkomnen/iot-bridge/cmd/valkyrie/device"
	"github.com/dkomnen/iot-bridge/proto"
	"github.com/dkomnen/iot-bridge/utils/constants"
	"github.com/golang/protobuf/proto"
)

type Thermometer struct {
	opts      device.Options
	isRunning bool
	stop      chan struct{}
}

func (t *Thermometer) Setup() error {
	if err := t.opts.Broker.Connect(); err != nil {
		return err
	}
	t.opts.Broker.Subscribe(
		"online",
		func(msg []byte) error {
			fmt.Println(string(msg))
			return nil
		},
	)
	t.sendOnlineMessage()
	return nil
}

func (t *Thermometer) Run() error {
	t.isRunning = true
	defer func() { t.isRunning = false }()

	tick := time.NewTicker(t.opts.Interval)
	for {
		select {
		case <-tick.C:
			t.generateMessage()
		case <-t.stop:
			return nil
		}
	}
	return nil
}

func (t *Thermometer) generateMessage() error {
	deviceReading := &message.DeviceReading{
		SensorType:    "temperature",
		Timestamp:     time.Now().Unix(),
		SerialNumber:  "1111",
		SensorReading: randomFloat32InRange(18, 22),
		Unit:          "c",
	}
	data, err := proto.Marshal(deviceReading)
	if err != nil {
		return err
	}

	t.opts.Broker.Publish(constants.ThermometerReading, data)

	return nil
}

func (t *Thermometer) sendOnlineMessage() error {
	deviceStatus := &message.DeviceStatus{
		DeviceStatus: true,
		SerialNumber: t.opts.SerialNumber,
	}
	data, err := proto.Marshal(deviceStatus)
	if err != nil {
		return err
	}
	t.opts.Broker.Publish(constants.DeviceStatus, data)

	return nil
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
		SerialNumber: "",
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
