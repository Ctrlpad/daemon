package ble

import (
	"github.com/charmbracelet/log"

	"tinygo.org/x/bluetooth"
)

func ScanAndConnectToCtrlpad() (*bluetooth.Device, error) {
	err := Adapter.Enable()
	if err != nil {
		log.Error("Adapter", "err", err)
	}
	log.Info("Enabled Adapter")
	deviceChan := make(chan bluetooth.ScanResult, 1)

	log.Info("Scanning for BLE device")
	err = Adapter.Scan(func(adapter *bluetooth.Adapter, result bluetooth.ScanResult) {
		log.Info("Found device", "Device Name", result.LocalName(), "RSSI", result.RSSI, "Address", result.Address.String())
		if result.LocalName() == "ctrlpad_BLE" {
			adapter.StopScan()
			deviceChan <- result
		}
	})
	if err != nil {
		return nil, err
	}

	foundDevice := <-deviceChan

	device, err := Adapter.Connect(foundDevice.Address, bluetooth.ConnectionParams{})
	if err != nil {
		return nil, err
	}
	log.Infof("Connected to %s", foundDevice.LocalName())
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
