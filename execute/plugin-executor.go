package execute

import (
	"os/exec"

	"github.com/charmbracelet/log"
)

type PluginExec struct {
	Path string
}

func NewPluginExecutor(pluginPath string) PluginExec {
	return PluginExec{
		Path: pluginPath,
	}
}

func (p PluginExec) Version() error {
	cmd := exec.Command(p.Path, "version")
	bts, err := cmd.Output()
	if err != nil {
		return err
	}
	log.Info(string(bts))
	return nil
}

func (p PluginExec) ConfigInit() error {
	return nil
}

func (p PluginExec) RunCommand(pluginName string, command string) error {
	cmd := exec.Command(pluginName, command)
	bts, err := cmd.Output()
	if err != nil {
		return err
	}
	log.Info(string(bts))
	return nil
}
