package execute

import (
	"fmt"
	"os/exec"

	"github.com/charmbracelet/log"
)

type PluginExec struct {
	Root string
}

func NewPluginExecutor() PluginExec {
	return PluginExec{}
}

// Executes a command and an argument to a Vision plugin binary then processes any output.
func (p PluginExec) RunCommand(pluginName string, arg string) error {
	root := fmt.Sprintf("vision-plugin-%v-v1", pluginName)
	cmd := exec.Command(root, arg)
	bts, err := cmd.Output()
	if err != nil {
		return err
	}
	log.Info(string(bts))
	return nil
}
