package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dkomnen/iot-bridge/broker"
	"github.com/dkomnen/iot-bridge/broker/mqtt"
	"github.com/dkomnen/iot-bridge/utils/constants"
	"github.com/dkomnen/iot-bridge/proto"
	"github.com/urfave/cli"

	"github.com/globalsign/mgo"
	"github.com/golang/protobuf/proto"
	"github.com/globalsign/mgo/bson"
)

var connAddr string

type temp struct {
	SerialNumber string  `bson:"serial_number"`
	Temperature  float64 `bson:"temperature"`
	Unit         string  `bson:"unit"`
	Timestamp    int64   `bson:"timestamp"`
}

func main() {
	mongoSession, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer mongoSession.Close()

	mongoSession.SetMode(mgo.Monotonic, true)

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
			constants.ThermometerReading,
			func(msg []byte) error {
				deviceReading := &message.DeviceReading{}
				err := proto.Unmarshal(msg, deviceReading)
				log.Printf("Device reading: %v\n", deviceReading)
				//saveReadingToMongo(mongoSession, deviceReading)
				if err != nil {
					return err
				}

				return nil
			},
		)
		b.Subscribe(
			constants.DeviceStatus,
			func(msg []byte) error {
				deviceStatus := &message.DeviceStatus{}
				err := proto.Unmarshal(msg, deviceStatus)
				log.Printf("Device status: %v\n", deviceStatus)
				UpdateDeviceStatus(mongoSession, deviceStatus.SerialNumber, deviceStatus.DeviceStatus)
				if err != nil {
					return err
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
	var serialNumber [32]byte
	var unit byte
	var temperature float64
	var timestamp int64

	if err := decode(
		bytes.NewReader(msg),
		&serialNumber,
		&unit,
		&temperature,
		&timestamp,
	); err != nil {
		return temp{}, fmt.Errorf("could not parse message: %v", err)
	}

	return temp{
		SerialNumber: fmt.Sprintf("%x", serialNumber),
		Unit:         fmt.Sprintf("%c", unit),
		Temperature:  temperature,
		Timestamp:    timestamp,
	}, nil
}

func decode(r io.Reader, data ...interface{}) error {
	for _, v := range data {
		if err := binary.Read(r, binary.BigEndian, v); err != nil {
			return err
		}
	}
	return nil
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

func saveReadingToMongo(mongoSession *mgo.Session, data *message.DeviceReading) {
	c := mongoSession.DB("heimdall").C("test")

	index := mgo.Index{
		Key:        []string{"sensor_type", "timestamp", "sensor_reading", "unit", "serial_number"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err := c.EnsureIndex(index)
	if err != nil {
		panic(err)
	}
	c.Insert(bson.M{"serial_number": data.SerialNumber,
		"sensor_type": data.SensorType,
		"sensor_reading": data.SensorReading,
		"unit": data.Unit,
		"timestamp": data.Timestamp})
}

func UpdateDeviceStatus(mongoSession *mgo.Session, serialNumber string, status bool) {
	c := mongoSession.DB("heimdall").C("device_status")

	update := bson.M{"$set": bson.M{"serial_number": serialNumber, "status": status}}
	selector := bson.M{"serial_number": serialNumber}

	_, err := c.Upsert(selector, update)
	if err != nil {
		panic(err)
	}
}
