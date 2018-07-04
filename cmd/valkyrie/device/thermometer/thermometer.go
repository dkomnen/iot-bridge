package thermometer

import (
	"context"
	"time"
	"fmt"

	"github.com/dkomnen/iot-bridge/broker/mqtt"
	"github.com/dkomnen/iot-bridge/cmd/valkyrie/device"
	"github.com/dkomnen/iot-bridge/proto"
	"github.com/dkomnen/iot-bridge/utils/constants"
	"github.com/golang/protobuf/proto"
	"os"
)

type Thermometer struct {
	opts      device.Options
	isRunning bool
	stop      chan struct{}
	run       chan struct{}
}

func (t *Thermometer) Setup() error {
	if err := t.opts.Broker.Connect(); err != nil {
		return err
	}
	fmt.Println("Subscribed to " + t.opts.SerialNumber + constants.RemoteShutdown)
	t.opts.Broker.Subscribe(
		t.opts.SerialNumber+constants.RemoteShutdown,
		func(msg []byte) error {
			fmt.Println("Device received remote shutdown command from server")
			fmt.Println("Shutting down...")
			t.sendOfflineMessage()
			t.isRunning = false
			return nil
		},
	)
	t.opts.Broker.Subscribe(
		t.opts.SerialNumber+constants.RemotePowerOn,
		func(msg []byte) error {
			fmt.Println("Device received remote power on command from server")
			fmt.Println("Starting device...")
			t.sendOnlineMessage()
			t.isRunning = true
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
		if t.isRunning == true {
			select {
			case <-tick.C:
				t.generateMessage()
			case <-t.stop:
				fmt.Println("Device stopped...")
				t.isRunning = false
			}
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

func (t *Thermometer) sendOfflineMessage() error {
	deviceStatus := &message.DeviceStatus{
		DeviceStatus: false,
		SerialNumber: t.opts.SerialNumber,
	}
	data, err := proto.Marshal(deviceStatus)
	if err != nil {
		return err
	}
	t.opts.Broker.Publish(constants.DeviceStatus, data)

	return nil
}

func (t *Thermometer) Stop() error {
	t.sendOfflineMessage()
	os.Exit(0)

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
