package device

import (
	"context"
	"crypto/sha512"
	"time"

	"github.com/dkomnen/iot-bridge/broker"
)

type Options struct {
	SerialNumber [32]byte
	Broker       broker.Broker
	Interval     time.Duration

	Custom context.Context
}

type Option func(*Options)

func SerialNumber(s [32]byte) Option {
	return func(opts *Options) {
		opts.SerialNumber = s
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

func GenerateSerialNumber(base string) [32]byte {
	return sha512.Sum512_256([]byte(base))
}
