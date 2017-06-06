package grovepi

import (
	"fmt"
	"time"

	"github.com/mrmorphic/hwio"
)

const (
	//Pins
	A0 = 0
	A1 = 1
	A2 = 2

	D2 = 2
	D3 = 3
	D4 = 4
	D5 = 5
	D6 = 6
	D7 = 7
	D8 = 8

	//Cmd format
	DIGITAL_READ  = 1
	DIGITAL_WRITE = 2
	ANALOG_READ   = 3
	ANALOG_WRITE  = 4
	PIN_MODE      = 5
	DHT_READ      = 40
)

type GrovePi struct {
	i2cmodule hwio.I2CModule
	i2cDevice hwio.I2CDevice
}

func InitGrovePi(address int) *GrovePi {
	grovePi := new(GrovePi)
	m, err := hwio.GetModule("i2c")
	if err != nil {
		fmt.Printf("could not get i2c module: %s\n", err)
		return nil
	}
	grovePi.i2cmodule = m.(hwio.I2CModule)
	grovePi.i2cmodule.Enable()

	grovePi.i2cDevice = grovePi.i2cmodule.GetDevice(address)
	return grovePi
}

func (grovePi *GrovePi) CloseDevice() {
	grovePi.i2cmodule.Disable()
}

func (grovePi *GrovePi) AnalogRead(pin byte) (int, error) {
	b := []byte{ANALOG_READ, pin, 0, 0}
	err := grovePi.i2cDevice.Write(1, b)
	if err != nil {
		return 0, err
	}
	time.Sleep(100 * time.Millisecond)
	grovePi.i2cDevice.ReadByte(1)
	val, err2 := grovePi.i2cDevice.Read(1, 4)
	if err2 != nil {
		return 0, err
	}
	var v1 int = int(val[1])
	var v2 int = int(val[2])
	return ((v1 * 256) + v2), nil
}

func (grovePi *GrovePi) DigitalRead(pin byte) (byte, error) {
	b := []byte{DIGITAL_READ, pin, 0, 0}
	err := grovePi.i2cDevice.Write(1, b)
	if err != nil {
		return 0, err
	}
	time.Sleep(100 * time.Millisecond)
	val, err2 := grovePi.i2cDevice.ReadByte(1)
	if err2 != nil {
		return 0, err2
	}
	return val, nil
}

func (grovePi *GrovePi) DigitalWrite(pin byte, val byte) error {
	b := []byte{DIGITAL_WRITE, pin, val, 0}
	err := grovePi.i2cDevice.Write(1, b)
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		return err
	}
	return nil
}

func (grovePi *GrovePi) PinMode(pin byte, mode string) error {
	var b []byte
	if mode == "output" {
		b = []byte{PIN_MODE, pin, 1, 0}
	} else {
		b = []byte{PIN_MODE, pin, 0, 0}
	}
	err := grovePi.i2cDevice.Write(1, b)
	time.Sleep(100 * time.Millisecond)
	if err != nil {
		return err
	}
	return nil
}

func (grovePi *GrovePi) ReadDHT(pin byte) ([]byte, error) {
	b := []byte{DHT_READ, pin, 0, 0}
	return grovePi.readDHTRawData(b)
}

func (grovePi *GrovePi) readDHTRawData(cmd []byte) ([]byte, error) {

	err := grovePi.i2cDevice.Write(1, cmd)
	if err != nil {
		return nil, err
	}

	time.Sleep(600 * time.Millisecond)

	grovePi.i2cDevice.ReadByte(1)

	time.Sleep(100 * time.Millisecond)

	raw, err := grovePi.i2cDevice.Read(1, 9)
	if err != nil {
		return nil, err
	}
	return raw, nil
}
