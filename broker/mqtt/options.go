package mqtt

import (
	"context"
	"time"

	"github.com/dkomnen/iot-bridge/broker"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

const (
	autoReconnect = iota
	binWillTopic
	binWillPayload
	binWillQoS
	binWillRetained
	binWillProvided
	willTopic
	willPayload
	willQoS
	willRetained
	willProvided
	unsetWill
	cleanSession
	clientID
	connectTimeout
	connectLostHandler
	defaultPublishHandler
	keepAlive
	maxReconnectInterval
	messageChannelDepth
	onConnectHandler
	orderMatters
	password
	pingTimeout
	protocolVersion
	store
	username
	writeTimeout

	publishQoS
	publishRetained

	subscribeQoS
)

func AutoReconnect(a bool) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, autoReconnect, a)
	}
}

func BinaryWill(topic string, payload []byte, qos byte, retained bool) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, binWillProvided, true)
		opts.Custom = context.WithValue(opts.Custom, binWillTopic, topic)
		opts.Custom = context.WithValue(opts.Custom, binWillPayload, payload)
		opts.Custom = context.WithValue(opts.Custom, binWillQoS, qos)
		opts.Custom = context.WithValue(opts.Custom, binWillRetained, retained)
	}
}

func CleanSession(clean bool) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, cleanSession, clean)
	}
}

func ClientID(id string) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, clientID, id)
	}
}

func ConnectTimeout(t time.Duration) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, connectTimeout, t)
	}
}

func ConnectLostHandler(f mqtt.ConnectionLostHandler) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, connectLostHandler, f)
	}
}

func DefaultPublishHandler(f mqtt.MessageHandler) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, defaultPublishHandler, f)
	}
}

func KeepAlive(t time.Duration) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, keepAlive, t)
	}
}

func MaxReconnectInterval(t time.Duration) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, maxReconnectInterval, t)
	}
}

func MessageChannelDepth(d uint) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, messageChannelDepth, d)
	}
}

func OnConnectHandler(onConn mqtt.OnConnectHandler) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, onConnectHandler, onConn)
	}
}

func OrderMatters(b bool) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, orderMatters, b)
	}
}

func Password(pw string) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, password, pw)
	}
}

func PingTimeout(t time.Duration) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, pingTimeout, t)
	}
}

func ProtocolVersion(pv uint) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, protocolVersion, pv)
	}
}

func Store(s mqtt.Store) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, store, s)
	}
}

//func (o *ClientOptions) SetUsername(u string) *ClientOptions
func Username(u string) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, username, u)
	}
}

func Will(topic string, payload string, qos byte, retained bool) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, willProvided, true)
		opts.Custom = context.WithValue(opts.Custom, willTopic, topic)
		opts.Custom = context.WithValue(opts.Custom, willPayload, payload)
		opts.Custom = context.WithValue(opts.Custom, willQoS, qos)
		opts.Custom = context.WithValue(opts.Custom, willRetained, retained)
	}
}

func WriteTimeout(t time.Duration) broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, writeTimeout, t)
	}
}

func UnsetWill() broker.Option {
	return func(opts *broker.Options) {
		opts.Custom = context.WithValue(opts.Custom, unsetWill, true)
	}
}

func PublishQoS(qos byte) broker.PublishOption {
	return func(opts *broker.PublishOptions) {
		opts.Custom = context.WithValue(opts.Custom, publishQoS, qos)
	}
}

func PublishRetained(b bool) broker.PublishOption {
	return func(opts *broker.PublishOptions) {
		opts.Custom = context.WithValue(opts.Custom, publishRetained, b)
	}
}

func SubscribeQoS(qos byte) broker.SubscribeOption {
	return func(opts *broker.SubscribeOptions) {
		opts.Custom = context.WithValue(opts.Custom, subscribeQoS, qos)
	}
}

func setBinaryWill(co *mqtt.ClientOptions, ctx context.Context) {
	var topic string
	var payload []byte
	var qos byte
	var retained, ok, okToSet bool

	okToSet = true
	if topic, ok = ctx.Value(binWillTopic).(string); !ok {
		okToSet = false
	}
	if payload, ok = ctx.Value(binWillPayload).([]byte); !ok {
		okToSet = false
	}
	if qos, ok = ctx.Value(binWillQoS).(byte); !ok {
		okToSet = false
	}
	if retained, ok = ctx.Value(binWillRetained).(bool); !ok {
		okToSet = false
	}

	if okToSet {
		co.SetBinaryWill(topic, payload, qos, retained)
	}
}

func setWill(co *mqtt.ClientOptions, ctx context.Context) {
	if v := ctx.Value(willProvided); v != nil {
		var topic string
		var payload string
		var qos byte
		var retained, ok, okToSet bool

		okToSet = true
		if topic, ok = ctx.Value(willTopic).(string); !ok {
			okToSet = false
		}
		if payload, ok = ctx.Value(willPayload).(string); !ok {
			okToSet = false
		}
		if qos, ok = ctx.Value(willQoS).(byte); !ok {
			okToSet = false
		}
		if retained, ok = ctx.Value(willRetained).(bool); !ok {
			okToSet = false
		}

		if okToSet {
			co.SetWill(topic, payload, qos, retained)
		}
	}
}

func setAutoReconnect(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(autoReconnect).(bool); ok {
		co.SetAutoReconnect(v)
	}
}

func setCleanSession(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(cleanSession).(bool); ok {
		co.SetCleanSession(v)
	}
}

func setClientID(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(clientID).(string); ok {
		co.SetClientID(v)
	}
}

func setConnectTimeout(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(connectTimeout).(time.Duration); ok {
		co.SetConnectTimeout(v)
	}
}

func setConnectLostHandler(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(connectLostHandler).(mqtt.ConnectionLostHandler); ok {
		co.SetConnectionLostHandler(v)
	}
}

func setDefaultPublishHandler(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(defaultPublishHandler).(mqtt.MessageHandler); ok {
		co.SetDefaultPublishHandler(v)
	}
}

func setKeepAlive(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(keepAlive).(time.Duration); ok {
		co.SetKeepAlive(v)
	}
}

func setMaxReconnectInterval(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(maxReconnectInterval).(time.Duration); ok {
		co.SetMaxReconnectInterval(v)
	}
}

func setMessageChannelDepth(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(messageChannelDepth).(uint); ok {
		co.SetMessageChannelDepth(v)
	}
}

func setOnConnectHandler(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(onConnectHandler).(mqtt.OnConnectHandler); ok {
		co.SetOnConnectHandler(v)
	}
}

func setOrderMatters(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(orderMatters).(bool); ok {
		co.SetOrderMatters(v)
	}
}

func setPassword(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(password).(string); ok {
		co.SetPassword(v)
	}
}

func setPingTimeout(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(pingTimeout).(time.Duration); ok {
		co.SetPingTimeout(v)
	}
}

func setProtocolVersion(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(protocolVersion).(uint); ok {
		co.SetProtocolVersion(v)
	}
}

func setStore(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(store).(mqtt.Store); ok {
		co.SetStore(v)
	}
}

func setUsername(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(username).(string); ok {
		co.SetUsername(v)
	}
}

func setWriteTimeout(co *mqtt.ClientOptions, ctx context.Context) {
	if v, ok := ctx.Value(writeTimeout).(time.Duration); ok {
		co.SetWriteTimeout(v)
	}
}

func setUnsetWill(co *mqtt.ClientOptions, ctx context.Context) {
	if _, ok := ctx.Value(unsetWill).(bool); ok {
		co.UnsetWill()
	}
}

func setAllCustomClientOptions(co *mqtt.ClientOptions, ctx context.Context) {
	opts := []func(*mqtt.ClientOptions, context.Context){
		setBinaryWill,
		setWill,
		setAutoReconnect,
		setCleanSession,
		setClientID,
		setConnectTimeout,
		setConnectLostHandler,
		setDefaultPublishHandler,
		setKeepAlive,
		setMaxReconnectInterval,
		setMessageChannelDepth,
		setOnConnectHandler,
		setOrderMatters,
		setPassword,
		setPingTimeout,
		setProtocolVersion,
		setStore,
		setUsername,
		setWriteTimeout,
		setUnsetWill,
	}

	for _, o := range opts {
		o(co, ctx)
	}
}
