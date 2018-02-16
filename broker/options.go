package broker

import "crypto/tls"

type Options struct {
	Address   string
	Secure    bool
	TLSConfig *tls.Config

	// custom options for the broker
	Custom map[string]interface{}
}

type Option func(*Options)

func Address(address string) Option {
	return func(opts *Options) {
		opts.Address = address
	}
}

func Secure(secure bool) Option {
	return func(opts *Options) {
		opts.Secure = secure
	}
}

func TLSConfig(t *tls.Config) Option {
	return func(opts *Options) {
		opts.TLSConfig = t
	}
}

type PublishOption func(*PublishOptions)

type PublishOptions struct {
	Custom map[string]interface{}
}

type SubscribeOption func(*SubscribeOptions)

type SubscribeOptions struct {
	Custom map[string]interface{}
}
