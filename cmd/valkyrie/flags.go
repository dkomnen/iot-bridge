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

var brokerUsernameFlag = cli.StringFlag{
	Name:  "broker-username",
	Usage: "set the username for the broker",
	Value: "",
}

var brokerPasswordFlag = cli.StringFlag{
	Name:  "broker-password",
	Usage: "set the password for the broker",
	Value: "",
}

var serialNumberFlag = cli.StringFlag{
	Name:  "serial-number",
	Usage: "string which will be hashed and used as serial number of the device",
	Value: "",
}
