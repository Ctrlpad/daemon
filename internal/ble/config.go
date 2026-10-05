package ble

import "tinygo.org/x/bluetooth"

var (
	Adapter = bluetooth.DefaultAdapter

	CtrlpadServiceUUID, _        = bluetooth.ParseUUID("a3308e24-786f-40b3-bf31-308875404027")
	CtrlpadCharacteristicUUID, _ = bluetooth.ParseUUID("62148466-62a9-4f65-bc29-2c2e408b8684")
)
