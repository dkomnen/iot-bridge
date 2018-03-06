package main

import (
	"github.com/dkomnen/iot-bridge/broker"
	"github.com/dkomnen/iot-bridge/cmd/valkyrie/device"
	"github.com/dkomnen/iot-bridge/cmd/valkyrie/device/thermometer"
	"github.com/dkomnen/iot-bridge/mqtt"
	"github.com/urfave/cli"
)

var (
	lowerLimitFlag = cli.Float64Flag{
		Name:  "lower-limit",
		Usage: "lower limit for `temperature` data that is sent",
		Value: 0.0,
	}
	upperLimitFlag = cli.Float64Flag{
		Name:  "upper-limit",
		Usage: "upper limit for `temperature` data that is sent",
		Value: 100.0,
	}
	fahrenheitFlag = cli.BoolFlag{
		Name:  "fahrenheit",
		Usage: "if set, the unit of measurement will be fahrenheit",
	}
	thermometerCommand = cli.Command{
		Name:        "thermometer",
		Usage:       "runs a simulation of a smart thermometer",
		Description: "simulates a thermometer that publishes data in a given interval",
		Flags: []cli.Flag{
			lowerLimitFlag,
			upperLimitFlag,
			fahrenheitFlag,
			intervalFlag,
			brokerAddressFlag,
			serialNumberFlag,
		},
		Action: runThermometer,
	}
)

func runThermometer(ctx *cli.Context) error {
	therm := makeThermometer(ctx)

	if err := run(therm); err != nil {
		return err
	}

	return nil

}

func makeThermometer(ctx *cli.Context) device.Device {
	brokerAddress := ctx.String(brokerAddressFlag.Name)
	b := mqtt.New(broker.Address(brokerAddress))

	return thermometer.New(
		device.Broker(b),
		device.SerialNumber(
			device.GenerateSerialNumber(ctx.String(serialNumberFlag.Name)),
		),
		device.Interval(ctx.Duration(intervalFlag.Name)),
		thermometer.DataRange(
			ctx.Float64(lowerLimitFlag.Name),
			ctx.Float64(upperLimitFlag.Name),
		),
		thermometer.Fahrenheit(ctx.Bool(fahrenheitFlag.Name)),
	)
}
