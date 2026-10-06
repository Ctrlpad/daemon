package ble

import (
	"github.com/charmbracelet/log"

	"tinygo.org/x/bluetooth"
)

func ConnectToCtrlpad(targetDevice string) (*bluetooth.Device, error) {
	err := Adapter.Enable()
	if err != nil {
		log.Error("Adapter", "err", err)
	}
	log.Info("Enabled Adapter")

	mac, err := bluetooth.ParseMAC(targetDevice)
	if err != nil {
		log.Error("Parsing mac", "err", err)
		return nil, err
	}

	addr := bluetooth.Address{
		MACAddress: bluetooth.MACAddress{
			MAC: mac,
		},
	}

	log.Info("Connecting to device", "Address", addr)
	device, err := Adapter.Connect(addr, bluetooth.ConnectionParams{})
	if err != nil {
		return nil, err
	}
	log.Infof("Connected to %s", targetDevice)

	return &device, nil
}

func SetupNotifications(device *bluetooth.Device) (chan string, error) {
	srvcs, err := device.DiscoverServices([]bluetooth.UUID{CtrlpadServiceUUID})
	if err != nil {
		return nil, err
	}
	srvc := srvcs[0]
	log.Info("Found service", "UUID", srvc.UUID().String())

	chars, err := srvc.DiscoverCharacteristics([]bluetooth.UUID{CtrlpadCharacteristicUUID})
	if err != nil {
		return nil, err
	}
	char := chars[0]
	log.Info("Found characteristic", "UUID", char.UUID().String())

	notifyChan := make(chan string, 1)

	err = char.EnableNotifications(func(buf []byte) {
		notifyChan <- string(buf)
	})
	if err != nil {
		return nil, err
	}
	log.Info("BLE Notifications enabled")

	return notifyChan, nil
}
