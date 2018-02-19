package mqtt

import (
	"context"

	"github.com/dkomnen/iot-bridge/broker"
)

const (
	autoReconnect = "autoReconnect"
)

//func (o *ClientOptions) SetAutoReconnect(a bool) *ClientOptions
func AutoReconnect(a bool) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, autoReconnect, a)
	}
}

//func (o *ClientOptions) SetBinaryWill(topic string, payload []byte, qos byte, retained bool) *ClientOptions
//func (o *ClientOptions) SetCleanSession(clean bool) *ClientOptions
//func (o *ClientOptions) SetClientID(id string) *ClientOptions
//func (o *ClientOptions) SetConnectTimeout(t time.Duration) *ClientOptions
//func (o *ClientOptions) SetConnectionLostHandler(onLost ConnectionLostHandler) *ClientOptions
//func (o *ClientOptions) SetDefaultPublishHandler(defaultHandler MessageHandler) *ClientOptions
//func (o *ClientOptions) SetKeepAlive(k time.Duration) *ClientOptions
//func (o *ClientOptions) SetMaxReconnectInterval(t time.Duration) *ClientOptions
//func (o *ClientOptions) SetMessageChannelDepth(s uint) *ClientOptions
//func (o *ClientOptions) SetOnConnectHandler(onConn OnConnectHandler) *ClientOptions
//func (o *ClientOptions) SetOrderMatters(order bool) *ClientOptions
//func (o *ClientOptions) SetPassword(p string) *ClientOptions
//func (o *ClientOptions) SetPingTimeout(k time.Duration) *ClientOptions
//func (o *ClientOptions) SetProtocolVersion(pv uint) *ClientOptions
//func (o *ClientOptions) SetStore(s Store) *ClientOptions
//func (o *ClientOptions) SetUsername(u string) *ClientOptions
//func (o *ClientOptions) SetWill(topic string, payload string, qos byte, retained bool) *ClientOptions
//func (o *ClientOptions) SetWriteTimeout(t time.Duration) *ClientOptions
//func (o *ClientOptions) UnsetWill() *ClientOptions
