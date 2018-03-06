package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

var (
	app = cli.NewApp()
)

func init() {
	app.Name = "Valkyrie"
	app.HelpName = "valkyrie"
	app.Description = "Valkyrie simulates various types of IOT devices"
	app.Usage = "mock an IOT device for testing"
	app.Version = "0.0.0"
	app.Author = "David Komljenovic"

	app.Commands = cli.Commands{
		thermometerCommand,
	}
}

func main() {
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
