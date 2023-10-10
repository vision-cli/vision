package internal

import (
	"fmt"
	"os/exec"
)

type Executor struct {
}

func (exe *Executor) RunCommand(updateStr string) error {
	cmd := exec.Command("curl", updateStr)
	_, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("error updating: %w", err)
	}

}

// method 1: run curl command for releases different to current version
// method 2: run go get @latest
// method 3: run curl command for latest
// deciding on all three: method 2 first, if that fails, method 1, then method 3
// curl -L https://github.com/charmbracelet/log/archive/refs/tags/v0.2.4.tar.gz > charmbracelet-v0.2.4.tar.gz
