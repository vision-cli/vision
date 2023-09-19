package comms_test

import (
	"fmt"
	"github.com/vision-cli/vision/common/execute"
	"github.com/vision-cli/vision/common/tmpl"
	"testing"

	"github.com/stretchr/testify/require"
	api_v1 "github.com/vision-cli/api/v1"
	"github.com/vision-cli/vision/common/comms"
	"github.com/vision-cli/vision/common/mocks"
	"github.com/vision-cli/vision/common/plugins"
)

var plugin = plugins.Plugin{
	Name:            "plugin",
	PluginPath:      "plugin",
	InternalCommand: nil,
}

func TestCall_WhenRequestIsNil_ReturnsError(t *testing.T) {
	_, err := comms.Call[int](plugin, nil, nil)
	require.Error(t, err)
}

func TestCall_WhenExecutorIsNil_ReturnsError(t *testing.T) {
	_, err := comms.Call[int](plugin, &pluginRequest, nil)
	require.Error(t, err)
}

func TestCall_WhenReturnIsValid_ReturnsStruct(t *testing.T) {
	e := mocks.NewMockExecutor()
	e.SetOutput(`{"Msg":"hello"}`)
	result, err := comms.Call[TestMsg](plugin, &pluginRequest, &e)
	require.NoError(t, err)
	require.Equal(t, "hello", result.Msg)
}

func TestCall_WhenReturnIsInvalid_ReturnsError(t *testing.T) {
	e := mocks.NewMockExecutor()
	e.SetOutput(`{"Result":"","Error":"some error"}`)
	_, err := comms.Call[TestMsg](plugin, &pluginRequest, &e)
	require.Error(t, err)
	require.Equal(t, "some error", err.Error())
}

func TestCall_WhenReturnIsNotExpectStructButAlsoNotError_ReturnsError(t *testing.T) {
	e := mocks.NewMockExecutor()
	e.SetOutput(`{"Result":"some result","Error":""}`)
	_, err := comms.Call[TestMsg](plugin, &pluginRequest, &e)
	require.Error(t, err)
	require.Equal(t, "did not get expected result type from plugin:, got PluginResponse with result some result", err.Error())
}

func TestCall_WhenInternalReturnIsValid_ReturnsStruct(t *testing.T) {
	e := mocks.NewMockExecutor()
	plugin := plugins.Plugin{
		Name:       "internal-plugin",
		PluginPath: "internal-plugin",
		InternalCommand: func(_ string, _ execute.Executor, _ tmpl.TmplWriter) string {
			return `{"Msg":"hello"}`
		},
	}
	result, err := comms.Call[TestMsg](plugin, &pluginRequest, &e)
	require.NoError(t, err)
	require.Equal(t, "hello", result.Msg)
}

func TestCall_WhenExecutorFails_ReturnsError(t *testing.T) {
	e := mocks.NewMockExecutor()
	e.SetOutputErr(fmt.Errorf("error"))
	_, err := comms.Call[int](plugin, &pluginRequest, &e)
	require.Error(t, err)
}

var pluginRequest = api_v1.PluginRequest{
	Command: "run",
	Args:    []string{},
}

type TestMsg struct {
	Msg string
}
