package device

import (
	"context"
	"crypto/rand"
	"time"

	"github.com/dkomnen/iot-bridge/broker"
	"fmt"
)

type Options struct {
	SerialNumber string
	Broker       broker.Broker
	Interval     time.Duration

	Custom context.Context
}

type Option func(*Options)

func SerialNumber(s string) Option {
	if s != "" {
		return func(opts *Options) {
			opts.SerialNumber = string(s)
		}
	}
	return func(opts *Options) {
		opts.SerialNumber = GenerateSerialNumber()
	}
}

func Broker(b broker.Broker) Option {
	return func(opts *Options) {
		opts.Broker = b
	}
}

func Interval(d time.Duration) Option {
	return func(opts *Options) {
		opts.Interval = d
	}
}

func GenerateSerialNumber() string {
	//return sha512.Sum512_256([]byte(base))
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		fmt.Println("error:", err)
		return ""
	}
	// Set the local bit
	buf[0] |= 2
	fmt.Printf("Generated MAC address: %02x:%02x:%02x:%02x:%02x:%02x\n", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
	return fmt.Sprintf("%02x:%02x:%02x:%02x:%02x:%02x", buf[0], buf[1], buf[2], buf[3], buf[4], buf[5])
}
