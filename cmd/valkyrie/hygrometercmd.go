package main

import (
"github.com/dkomnen/iot-bridge/broker"
"github.com/dkomnen/iot-bridge/broker/mqtt"
"github.com/dkomnen/iot-bridge/cmd/valkyrie/device"
"github.com/dkomnen/iot-bridge/cmd/valkyrie/device/hygrometer"
"github.com/dkomnen/iot-bridge/utils/constants"

"github.com/urfave/cli"
)

var (
	hygrometerCommand = cli.Command{
		Name:        "hygrometer",
		Usage:       "runs a simulation of a smart hygrometer",
		Description: "simulates a hygrometer that publishes data in a given interval",
		Flags: []cli.Flag{
			lowerLimitFlag,
			upperLimitFlag,
			intervalFlag,
			brokerAddressFlag,
			serialNumberFlag,
		},
		Action: runHygrometer,
	}
)

func runHygrometer(ctx *cli.Context) error {
	therm := makeHygrometer(ctx)

	if err := run(therm); err != nil {
		return err
	}

	return nil

}

func makeHygrometer(ctx *cli.Context) device.Device {
	brokerAddress := ctx.String(brokerAddressFlag.Name)
	brokerUsername := ctx.String(brokerUsernameFlag.Name)
	brokerPassword := ctx.String(brokerPasswordFlag.Name)

	lastWillMessage := getLastWillMessage(ctx.String(serialNumberFlag.Name))

	b := mqtt.New(broker.Address(brokerAddress), mqtt.Username(brokerUsername), mqtt.Password(brokerPassword), mqtt.BinaryWill(constants.DeviceStatus, lastWillMessage, 1, false))
	return hygrometer.New(
		device.Broker(b),
		device.SerialNumber(
			ctx.String(serialNumberFlag.Name),
		),
		device.Interval(ctx.Duration(intervalFlag.Name)),
		hygrometer.DataRange(
			ctx.Float64(lowerLimitFlag.Name),
			ctx.Float64(upperLimitFlag.Name),
		),
	)
}



