package cmd

import (
	"github.com/charmbracelet/log"
	"github.com/ctrlpad/daemon/internal/ble"
	"github.com/ctrlpad/daemon/internal/executor"
)

func Run() int {
	device, err := ble.ScanAndConnectToCtrlpad()
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
