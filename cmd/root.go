package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/ctrlpad/daemon/internal/ble"
	"github.com/ctrlpad/daemon/internal/executor"
	"github.com/ctrlpad/daemon/internal/info"
	"github.com/spf13/cobra"
)

var (
	targetDevice string
	verbose      bool
	rootCmd      = &cobra.Command{
		Use:   "ctrlpad-daemon",
		Short: "Daemon of the Ctrlpad. Executes the button actins (Execute Applicatin, change Volume, etc.)",
		Long:  `The daemon is the background service behind Ctrlpad. It scans for the Ctrlpad_BLE device, connects to it over Bluetooth Low Energy, subscribes to the button characteristic and executes the action that comes with every button press.`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if verbose {
				log.SetLevel(log.DebugLevel)
			}
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return Run(targetDevice)
		},
	}
)

func Execute() error {
	return rootCmd.Execute()
}

func Init(version, commit, date string) {
	cobra.OnInitialize(func() {
		log.SetDefault(log.NewWithOptions(os.Stderr, log.Options{
			ReportTimestamp: false,
			ReportCaller:    false,
			Level:           log.InfoLevel,
		}))
		info.SetVersionInfo(version, commit, date)
	})
}

func init() {
	rootCmd.Flags().StringVarP(&targetDevice, "device", "d", "", "Device to connect")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.MarkFlagRequired("device")
}

func Run(targetDevice string) error {
	device, err := ble.ConnectToCtrlpad(targetDevice)
	if err != nil {
		return err
	}

	buttonConfigs, err := ble.SetupNotifications(device)
	if err != nil {
		return err
	}

	for buttonConfig := range buttonConfigs {
		err := executor.ExecuteAction(buttonConfig)
		if err != nil {
			return err
		}
	}
	return err
}
