package device

type Device interface {
	Setup() error
	Run() error
	Stop() error
	Options() Options
}
