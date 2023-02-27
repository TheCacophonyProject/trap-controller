/*
trapdbusclient - client for accessing Cacophony events
Copyright (C) 2020, The Cacophony Project

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

package trapdbusclient

import (
	"encoding/json"

	"github.com/godbus/dbus"
)

func TriggerTrap(details map[string]interface{}) error {
	detailsBytes, err := json.Marshal(details)
	if err != nil {
		return err
	}
	_, err = trapDbusCall(
		"org.cacophony.TrapController.TriggerTrap",
		string(detailsBytes))
	return err
}

func SetTrapPower(power bool) error {
	_, err := trapDbusCall(
		"org.cacophony.TrapController.TriggerTrap",
		power)
	return err
}

func trapDbusCall(method string, params ...interface{}) ([]interface{}, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	obj := conn.Object("org.cacophony.TrapController", "/org/cacophony/TrapController")
	call := obj.Call(method, 0, params...)
	return call.Body, call.Err
}
