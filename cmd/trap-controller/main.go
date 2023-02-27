/*
trap-controller - Communicates with trap
Copyright (C) 2023, The Cacophony Project

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/

package main

import (
	"errors"
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/TheCacophonyProject/event-reporter/v3/eventclient"
	"github.com/TheCacophonyProject/go-config"
	"github.com/alexflint/go-arg"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

const (
	powerTrapPin   = "GPIO27"
	triggerTrapPin = "GPIO10"
	irLightsPin    = "GPIO18"
)

var version = "No version provided"

type Args struct {
	ConfigDir string `arg:"-c,--config" help:"configuration folder"`
}

func procArgs() Args {
	args := Args{
		ConfigDir: config.DefaultConfigDir,
	}
	arg.MustParse(&args)
	return args
}

func (Args) Version() string {
	return version
}

func main() {
	if err := runMain(); err != nil {
		log.Fatal(err)
	}
	// If no error then keep the background goroutines running.
	runtime.Goexit()
}

func runMain() error {
	_ = procArgs()
	log.SetFlags(0)

	log.Printf("running version: %s", version)

	log.Println("starting dbus service")
	if err := startService(); err != nil {
		return err
	}

	if err := setIRLightsPower(true); err != nil {
		return err
	}
	if err := setTrapPower(true); err != nil {
		return err
	}
	if err := setPin(triggerTrapPin, false); err != nil {
		return err
	}

	return nil
}

func triggerTrap(details map[string]interface{}) error {
	log.Println("triggering trap")
	host.Init()
	triggerTrap := gpioreg.ByName(triggerTrapPin)
	if triggerTrap == nil {
		return errors.New("failed to get pin to set trap power")
	}
	if err := triggerTrap.Out(gpio.High); err != nil {
		return fmt.Errorf("failed to write to trap power pin, %e", err)
	}
	go func(pin gpio.PinIO) {
		time.Sleep(3 * time.Second)
		pin.Out(gpio.Low)
	}(triggerTrap)

	return eventclient.AddEvent(eventclient.Event{
		Timestamp: time.Now(),
		Type:      "trap-triggered",
		Details:   details,
	})
}

func setTrapPower(power bool) error {
	if power {
		log.Println("powering on trap")
	} else {
		log.Println("powering off trap")
	}
	return setPin(powerTrapPin, power)
}

func setIRLightsPower(power bool) error {
	if power {
		log.Println("powering on IR lights")
	} else {
		log.Println("powering off IR lights")
	}
	return setPin(irLightsPin, power)
}

func setPin(pinStr string, power bool) error {
	host.Init()
	var gpioLevel gpio.Level
	if power {
		gpioLevel = gpio.High
	} else {
		gpioLevel = gpio.Low
	}
	pin := gpioreg.ByName(pinStr)
	if pin == nil {
		return fmt.Errorf("failed to get pin '%s'", pinStr)
	}
	if err := pin.Out(gpioLevel); err != nil {
		return fmt.Errorf("failed write to pin '%s'", pinStr)
	}
	return nil

}
