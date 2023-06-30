package plugins

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

	api_v1 "github.com/vision-cli/api/v1"
	"github.com/vision-cli/vision/execute"
)

func Unmarshal[T any](in []byte) (*T, error) {
	var out T
	err := json.Unmarshal(in, &out)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func Marshal[T any](in *T) (string, error) {
	out, err := json.Marshal(in)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func Call[T any](plugin string, request *api_v1.PluginRequest, executor execute.Executor) (*T, error) {
	cmd := exec.Command(plugin)
	query, err := Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("cannot marshal request for plugin %s: %s", plugin, err.Error())
	}
	cmd.Stdin = strings.NewReader(query)
	response, err := executor.Output(cmd, ".", "calling plugin "+plugin)
	println(response)
	if err != nil {
		return nil, fmt.Errorf("cannot run plugin %s", plugin)
	}
	out, err := Unmarshal[T]([]byte(response))
	if err != nil {
		println(err.Error())
		// check if the response is an error
		outerr, err := Unmarshal[api_v1.PluginResponse]([]byte(response))
		if err != nil {
			return nil, fmt.Errorf("cannot unmarshal response from plugin %s: %s", plugin, err.Error())
		}
		if outerr.Error != "" {
			return nil, fmt.Errorf(outerr.Error)
		}
		return nil,
			fmt.Errorf("did not get expected result type from %s:, got PluginResponse with result %s", plugin, outerr.Result)
	}
	return out, nil
}
