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
		Name:        "sample-rate-ms",
		Type:        describe.Int32,
		Required:    false,
		Default:     int32(10000),
		Description: "Sample rate in milliseconds",
	}

	grovepi_pin := describe.Parameter{
		Name:        "pin",
		Type:        describe.String,
		Required:    true,
		Description: "Pin to read data from",
		Examples:    []string{"A0", "A1", "A2", "D2", "D3", "D4", "D5", "D6", "D7", "D8"},
	}

	hub.RegisterListener("grovepi-dht",
		func(config describe.Values) (hub.Listener, error) {

			hertz := config.Int32WithDefault(grovepi_samplerate.Name, grovepi_samplerate.Default.(int32))

			pinStr := config.MustString(grovepi_pin.Name)
			var pin byte
			switch pinStr {
			case "AO":
				pin = byte(A0)
			case "A1":
				pin = byte(A1)
			case "A2":
				pin = byte(A2)
			case "D2":
				pin = byte(D2)
			case "D3":
				pin = byte(D3)
			case "D4":
				pin = byte(D4)
			case "D5":
				pin = byte(D5)
			case "D6":
				pin = byte(D6)
			case "D7":
				pin = byte(D7)
			default:
				return nil, fmt.Errorf("unknown pin : %s", pinStr)
			}

			return newDHTListener(hertz, pin)
		},
		describe.Parameters{
			grovepi_samplerate,
			grovepi_pin,
		})

}

func newDHTListener(sampleTimeInMs int32, pin byte) (*dhtListener, error) {
	return &dhtListener{
		sampleTimeInMs: sampleTimeInMs,
		pin:            pin,
		close:          make(chan struct{}),
	}, nil
}

type dhtListener struct {
	sampleTimeInMs int32
	pin            byte
	close          chan struct{}
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
	wait := time.Millisecond * time.Duration(h.sampleTimeInMs)

	for {
		select {
		case <-h.close:
			return
		case <-time.Tick(wait):

		}

		data, err := grove.ReadDHT(h.pin)

		if err == nil {
			channel.Out() <- listener.NewHubMessage(data, "GROVEPI", "TODO : CHANGEME")
		} else {
			channel.Errors() <- err
		}

	}

}
