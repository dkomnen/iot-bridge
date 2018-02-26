package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dkomnen/iot-bridge/mqtt"
)

func main() {
	b := mqtt.New(mqtt.ClientID("bifrost"))
	if err := b.Connect(); err != nil {
		log.Fatal(err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	b.Subscribe(
		"topic",
		func(msg []byte) error {
			log.Printf("New message: %s\n", msg)
			return nil
		},
	)
	log.Println("Bifrost listening...")

	<-shutdown
	log.Println("Shutting down...")
}
