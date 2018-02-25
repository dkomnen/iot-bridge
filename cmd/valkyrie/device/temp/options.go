package temp

import (
	"context"

	"github.com/dkomnen/iot-bridge/cmd/valkyrie/device"
)

const (
	lowerBound = iota
	higherBound
	fahrenheit
)

func DataRange(low, high float64) device.Option {
	return func(opts *device.Options) {
		opts.Custom = context.WithValue(opts.Custom, lowerBound, low)
		opts.Custom = context.WithValue(opts.Custom, higherBound, high)
	}
}

func Fahrenheit(b bool) device.Option {
	return func(opts *device.Options) {
		opts.Custom = context.WithValue(opts.Custom, fahrenheit, b)
	}
}
