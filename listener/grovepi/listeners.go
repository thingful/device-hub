package grovepi

import (
	"errors"
	"fmt"
	"sync"
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

	hub.RegisterListener("grovepi-dht",
		func(config describe.Values) (hub.Listener, error) {

			hertz := config.Int32WithDefault(grovepi_samplerate.Name, grovepi_samplerate.Default.(int32))
			return newDHTListener(hertz)
		},
		describe.Parameters{
			grovepi_samplerate,
		})

}

func newDHTListener(sampleTimeInMs int32) (*dhtListener, error) {
	return &dhtListener{
		sampleTimeInMs: sampleTimeInMs,
		close:          make(chan struct{}),
	}, nil
}

type dhtListener struct {
	sampleTimeInMs int32
	close          chan struct{}
	started        bool
	lock           sync.Mutex
}

func (h *dhtListener) NewChannel(uri string) (hub.Channel, error) {

	h.lock.Lock()
	defer h.lock.Unlock()
	if h.started {
		return nil, errors.New("listener already monitoring an existing humidity sensor")
	}
	h.started = true

	var pin byte
	switch uri {
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
		return nil, fmt.Errorf("unknown uri : %s", uri)
	}

	errors := make(chan error)
	out := make(chan hub.Message)

	channel := listener.NewDefaultChannel(errors, out, func() error {
		return h.Close()
	})

	go h.loop(channel, pin)

	return channel, nil
}

func (h *dhtListener) Close() error {
	h.close <- struct{}{}
	return nil
}

func (h *dhtListener) loop(channel hub.Channel, pin byte) {

	// What is this magic number??
	grove := InitGrovePi(0x04)
	wait := time.Millisecond * time.Duration(h.sampleTimeInMs)

l:
	for {
		select {
		case <-h.close:
			break l
		case <-time.Tick(wait):

		}

		data, err := grove.ReadDHT(pin)

		if err == nil {
			channel.Out() <- listener.NewHubMessage(data, "GROVEPI", string(pin))
		} else {
			channel.Errors() <- err
		}

	}
}
