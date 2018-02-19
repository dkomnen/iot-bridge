package broker

import (
	"context"
	"crypto/tls"
)

type Options struct {
	Address   string
	TLSConfig *tls.Config

	// custom options for the broker
	Custom context.Context
}

type Option func(*Options)

func Address(address string) Option {
	return func(opts *Options) {
		opts.Address = address
	}
}

func TLSConfig(t *tls.Config) Option {
	return func(opts *Options) {
		opts.TLSConfig = t
	}
}

type PublishOption func(*PublishOptions)

type PublishOptions struct {
	Custom context.Context
}

type SubscribeOption func(*SubscribeOptions)

type SubscribeOptions struct {
	Custom context.Context
}
