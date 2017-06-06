package grovepi

import (
	"fmt"
	"time"

	hub "github.com/thingful/device-hub"
	"github.com/thingful/device-hub/describe"
	"github.com/thingful/device-hub/listener"
)

func init() {

	grovepi_samplerate := describe.Parameter{
		Name:        "grovepi-sample-rate",
		Type:        describe.Float32,
		Required:    false,
		Default:     1,
		Description: "Sample rate in Hertz e.g. times per second",
	}

	grovepi_pin := describe.Parameter{
		Name:        "grovepi-pin",
		Type:        describe.String,
		Required:    true,
		Description: "Pin to read data from",
		Examples:    []string{"A0", "A1", "A2", "D2", "D3", "D4", "D5", "D6", "D7", "D8"},
	}

	hub.RegisterListener("grovepi-dht",
		func(config describe.Values) (hub.Listener, error) {

			fmt.Println("grovepi-dht", config)

			hertz := 1
			pin := byte(D7)

			return newDHTListener(hertz, pin)
		},
		describe.Parameters{
			grovepi_samplerate,
			grovepi_pin,
		})

}

func newDHTListener(hertz int, pin byte) (*dhtListener, error) {
	return &dhtListener{
		hertz: hertz,
		pin:   pin,
		close: make(chan struct{}),
	}, nil
}

type dhtListener struct {
	hertz int
	pin   byte
	close chan struct{}
}

func (h *dhtListener) NewChannel(uri string) (hub.Channel, error) {

	errors := make(chan error)
	out := make(chan hub.Message)

	channel := listener.NewDefaultChannel(errors, out, func() error {
		return nil
	})

	go h.loop(channel)

	return channel, nil
}

func (h *dhtListener) Close() error {
	h.close <- struct{}{}
	return nil
}

func (h *dhtListener) loop(channel hub.Channel) {

	// What is this magic number??
	grove := InitGrovePi(0x04)

	for {
		select {
		case <-h.close:
			return
		case <-time.Tick(time.Millisecond * 1000):

		}

		data, err := grove.ReadDHT(h.pin)

		if err == nil {
			channel.Out() <- listener.NewHubMessage(data, "GROVEPI", "TODO : CHANGEME")
		} else {
			channel.Errors() <- err
		}

	}

}
