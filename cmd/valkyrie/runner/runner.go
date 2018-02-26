package runner

import (
	"fmt"

	"github.com/dkomnen/iot-bridge/cmd/valkyrie/device"
)

type Runner struct {
	devs   []device.Device
	Errors chan error
}

func (r *Runner) Run() {
	for i, dev := range r.devs {
		err := dev.Setup()
		if err != nil {
			r.Errors <- fmt.Errorf("not running device %d. device setup failed: %v", i, err)
		} else {
			go func(i int) {
				if err := dev.Run(); err != nil {
					r.Errors <- fmt.Errorf("error running device %d: %v", i, err)
				}
			}(i)
		}
	}
}

func (r *Runner) Stop() {
	for i, dev := range r.devs {
		if err := dev.Stop(); err != nil {
			r.Errors <- fmt.Errorf("error stopping device %d: %v", i, err)
		}
	}
	close(r.Errors)
}

func New(devs []device.Device) *Runner {
	return &Runner{devs: devs, Errors: make(chan error)}
}
