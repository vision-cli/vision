package plugins_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api_v1 "github.com/vision-cli/api/v1"
	"github.com/vision-cli/vision/plugins"
)

func TestGetCobraCommand_WithValidInputProvided_ReturnsObject(t *testing.T) {
	req := []byte(`{"Result":"result","Error":""}`)
	result, err := plugins.Unmarshal[api_v1.PluginResponse](req)
	expected := api_v1.PluginResponse{
		Result: "result",
		Error:  "",
	}
	require.NoError(t, err)
	assert.Equal(t, &expected, result)
}
