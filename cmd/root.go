package cmd

import (
	"flag"

	"github.com/charmbracelet/log"
	"github.com/ctrlpad/daemon/internal/ble"
	"github.com/ctrlpad/daemon/internal/executor"
)

func Run() int {
	var targetDevice string
	flag.StringVar(&targetDevice, "device", "", "Target device to connect to (MAC address)")
	flag.Parse()

	if targetDevice == "" {
		log.Error("Connection", "err", "Please pass an target device")
		return 1
	}

	device, err := ble.ConnectToCtrlpad(targetDevice)
	if err != nil {
		log.Error("Connection", "err", err)
		return 1
	}

	buttonConfigs, err := ble.SetupNotifications(device)
	if err != nil {
		log.Error("SetupNotifications", "err", err)
		return 1
	}

	for buttonConfig := range buttonConfigs {
		err := executor.ExecuteAction(buttonConfig)
		if err != nil {
			log.Error("Executor", "err", err)
		}
	}
	return 1
}
