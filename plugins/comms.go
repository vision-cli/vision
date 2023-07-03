package plugins

import (
	"fmt"
	"os/exec"
	"strings"

	api_v1 "github.com/vision-cli/api/v1"
	"github.com/vision-cli/common/execute"
	"github.com/vision-cli/common/marshal"
)

func Call[T any](plugin string, request *api_v1.PluginRequest, executor execute.Executor) (*T, error) {
	cmd := exec.Command(plugin)
	query, err := marshal.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request for plugin %s: %s", plugin, err.Error())
	}
	cmd.Stdin = strings.NewReader(query)
	response, err := executor.Output(cmd, ".", "calling plugin "+plugin)
	if err != nil {
		return nil, fmt.Errorf("cannot run plugin %s", plugin)
	}
	out, err := marshal.Unmarshal[T](response)
	if err != nil {
		// check if the response is an error
		outerr, err := marshal.Unmarshal[api_v1.PluginResponse](response)
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal response from plugin %s: %s", plugin, err.Error())
		}
		if outerr.Error != "" {
			return nil, fmt.Errorf(outerr.Error)
		}
		return nil,
			fmt.Errorf("did not get expected result type from %s:, got PluginResponse with result %s", plugin, outerr.Result)
	}
	return &out, nil
}
