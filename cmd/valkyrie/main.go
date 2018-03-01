package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dkomnen/iot-bridge/broker"
	"github.com/dkomnen/iot-bridge/cmd/valkyrie/device"
	"github.com/dkomnen/iot-bridge/cmd/valkyrie/device/temp"
	"github.com/dkomnen/iot-bridge/mqtt"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Valkyrie"
	app.HelpName = "valkyrie"
	app.Description = "Valkyrie simulates various types of IOT devices support"
	app.Usage = "mock an IOT device for testing Bifrost"
	app.Version = "0.0.0"
	app.Author = "David Komljenovic"

	app.Commands = cli.Commands{
		cli.Command{
			Name:        "temp",
			Aliases:     []string{"t"},
			Category:    "device simulator",
			Usage:       "runs a simulation of a smart thermometer",
			Description: "This command simulates a thermometer which sends data periodically to Bifrost in a given interval.",
			Flags: []cli.Flag{
				cli.DurationFlag{
					Name:  "tick, t",
					Usage: "specify the `interval` in milliseconds in which the messages will be sent",
					Value: time.Millisecond * 2500,
				},
				cli.Float64Flag{
					Name:  "lower-limit, l",
					Usage: "lower limit for `temperature` data that is sent",
					Value: 0.0,
				},
				cli.Float64Flag{
					Name:  "upper-limit, u",
					Usage: "upper limit for `temperature` data that is sent",
					Value: 100.0,
				},
				cli.UintFlag{
					Name:  "number-of-instances, n",
					Usage: "`number` of thermometer devices that run concurrently",
					Value: 1,
				},
				cli.BoolFlag{
					Name:  "fahrenheit, f",
					Usage: "if set, the unit of measurement will be fahrenheit",
				},
				cli.StringFlag{
					Name:  "broker-address",
					Usage: "set the address for the broker",
					Value: "tcp://localhost:1883",
				},
				cli.StringFlag{
					Name:  "serial-number, s",
					Usage: "string which will be hashed and used as serial number of the device",
					Value: "temp",
				},
			},
			Action: func(ctx *cli.Context) error {
				baddr := ctx.String("broker-address")
				log.Println("broker address:", baddr)
				b := mqtt.New(
					broker.Address(baddr),
					mqtt.ClientID("temp"),
				)

				t := temp.New(
					device.Broker(b),
					device.SerialNumber(
						device.GenerateSerialNumber(ctx.String("serial-number")),
					),
					device.Interval(ctx.Duration("tick")),
					temp.DataRange(
						ctx.Float64("lower-limit"),
						ctx.Float64("upper-limit"),
					),
					temp.Fahrenheit(ctx.Bool("fahrenheit")),
				)

				stop := make(chan os.Signal, 1)
				signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

				if err := t.Setup(); err != nil {
					return err
				}
				go func() {
					<-stop
					if err := t.Stop(); err != nil {
						log.Fatal(err)
					}
				}()

				if err := t.Run(); err != nil {
					return err
				}

				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
