package main

import (
	"time"

	"github.com/urfave/cli"
)

var intervalFlag = cli.DurationFlag{
	Name:  "interval",
	Usage: "specify the `duration` in milliseconds in which the messages will be sent",
	Value: time.Millisecond * 2500,
}

var brokerAddressFlag = cli.StringFlag{
	Name:  "broker-address",
	Usage: "set the address for the broker",
	Value: "tcp://localhost:1883",
}

var serialNumberFlag = cli.StringFlag{
	Name:  "serial-number",
	Usage: "string which will be hashed and used as serial number of the device",
	Value: "thermometer",
}
