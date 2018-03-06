package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/dkomnen/iot-bridge/cmd/valkyrie/device"
)

func run(dev device.Device) error {
	if err := dev.Setup(); err != nil {
		return err
	}

	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		defer signal.Stop(stop)

		<-stop
		dev.Stop()
	}()

	if err := dev.Run(); err != nil {
		return err
	}

	return nil
}
