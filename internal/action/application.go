package action

import (
	"os/exec"

	"github.com/ctrlpad/daemon/internal/desktop"
)

func ExecApplication(target string) error {
	desktop, err := desktop.Find(target)
	if err != nil {
		return err
	}
	cmd := exec.Command(desktop.Binary)
	err = cmd.Start()
	if err != nil {
		return err
	}
	go cmd.Wait()
	return nil
}
