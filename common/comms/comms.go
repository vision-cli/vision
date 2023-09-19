package comms

import (
	"fmt"
	"os/exec"
	"strings"

	api_v1 "github.com/vision-cli/api/v1"
	"github.com/vision-cli/vision/common/execute"
	"github.com/vision-cli/vision/common/marshal"
	"github.com/vision-cli/vision/common/plugins"
	"github.com/vision-cli/vision/common/tmpl"
)

func Call[T any](plugin plugins.Plugin, request *api_v1.PluginRequest, executor execute.Executor) (*T, error) {
	if request == nil {
		return nil, fmt.Errorf("comms.Call called with nil request")
	}
	if executor == nil {
		return nil, fmt.Errorf("comms.Call called with nil executor")
	}
	cmd := exec.Command(plugin.PluginPath)
	query, err := marshal.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request for plugin %s: %s", plugin.Name, err.Error())
	}
	cmd.Stdin = strings.NewReader(query)

	var response string
	if plugin.InternalCommand == nil {
		execResponse, err := executor.Output(cmd, ".", "calling plugin "+plugin.Name)
		if err != nil {
			return nil, fmt.Errorf("cannot run plugin %s", plugin.Name)
		}
		response = execResponse
	} else {
		templateWriter := tmpl.NewOsTmpWriter()
		response = plugin.InternalCommand(query, executor, templateWriter)
	}

	out, err := marshal.Unmarshal[T](response)
	if err != nil {
		// check if the response is an error
		outerr, err := marshal.Unmarshal[api_v1.PluginResponse](response)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal response from plugin %s: %s", plugin.Name, err.Error())
		}
		if outerr.Error != "" {
			return nil, fmt.Errorf(outerr.Error)
		}
		return nil,
			fmt.Errorf("did not get expected result type from %s:, got PluginResponse with result %s", plugin.Name, outerr.Result)
	}
	return &out, nil
}
