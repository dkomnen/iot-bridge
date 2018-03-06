package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/dkomnen/iot-bridge/broker"
	"github.com/dkomnen/iot-bridge/broker/mqtt"
	"github.com/urfave/cli"
)

var connAddr string

type temp struct {
	SerialNumber string  `json:"serial_number"`
	Temperature  float64 `json:"temperature"`
	Unit         string  `json:"unit"`
}

func main() {
	app := cli.NewApp()
	app.Name = "Bifrost"
	app.HelpName = "bifrost"
	app.Description = "Bifrost bridges the communication gap between various IOT devices and the connect server."
	app.Usage = "listen, parse and pass on all the messages"
	app.Version = "0.0.0"
	app.Author = "David Komljenovic"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "broker-address",
			Usage: "set the address of the broker",
			Value: "tcp://localhost:1883",
		},
		cli.StringFlag{
			Name:        "connect-address",
			Usage:       "set the address of the connect server",
			Value:       "http://localhost:8080",
			Destination: &connAddr,
		},
	}

	app.Action = func(ctx *cli.Context) error {
		b := mqtt.New(
			mqtt.ClientID("bifrost"),
			broker.Address(ctx.String("broker-address")),
		)
		if err := b.Connect(); err != nil {
			log.Fatal(err)
		}

		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

		b.Subscribe(
			"THERMOMETER",
			func(msg []byte) error {
				if t, err := parseTempMsg(msg); err != nil {
					return err
				} else {
					var buff bytes.Buffer
					if err := json.NewEncoder(&buff).Encode(&t); err != nil {
						return err
					}

					if err := sendToConnectAPI(connAddr, &buff); err != nil {
						return err
					}
				}

				return nil
			},
		)
		log.Println("Bifrost listening...")

		<-shutdown
		log.Println("Shutting down...")

		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func parseTempMsg(msg []byte) (temp, error) {
	var t temp
	t.SerialNumber = fmt.Sprintf("%x", msg[:32])
	if msg[32:33][0] == 'c' {
		t.Unit = "celsius"
	} else if msg[32:33][0] == 'f' {
		t.Unit = "fahrenheit"
	} else {
		return temp{}, fmt.Errorf(
			"malformed message: expected unit character to be 'c' or 'f', got %c",
			msg[32:33][0],
		)
	}
	if parsed, err := strconv.ParseFloat(string(msg[33:]), 64); err != nil {
		return temp{}, fmt.Errorf("malformed message: %v", err)
	} else {
		t.Temperature = parsed
	}

	return t, nil
}

func sendToConnectAPI(address string, msg io.Reader) error {
	req, err := http.NewRequest(http.MethodPost, address, msg)
	if err != nil {
		return err
	}

	req.Header.Set("content-type", "application/json")

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if rsp.StatusCode != http.StatusOK {
		return fmt.Errorf(
			"connect API responeded with %d: %s",
			rsp.StatusCode,
			rsp.Status,
		)
	}

	return nil
}
