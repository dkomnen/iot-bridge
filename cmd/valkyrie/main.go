package main

import (
	"log"
	"os"
	"time"

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

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "broker-address",
			Usage: "set the address for the broker",
			Value: "tcp://localhost:1883",
		},
	}

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
			},
			Action: func(ctx *cli.Context) error {
				// TODO: run the temp device type
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
