package plugins_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api_v1 "github.com/vision-cli/api/v1"
	"github.com/vision-cli/vision/plugins"
)

func TestUnmarshal_WithValidInputProvided_ReturnsObject(t *testing.T) {
	req := []byte(`{"Result":"result","Error":""}`)
	result, err := plugins.Unmarshal[api_v1.PluginResponse](req)
	require.NoError(t, err)
	expected := api_v1.PluginResponse{
		Result: "result",
		Error:  "",
	}
	assert.Equal(t, &expected, result)
}

func TestUnmarshal_ReturnsErrorWhenInValidInputProvided(t *testing.T) {
	req := []byte(`{"Result":"result","Error":"",}`) // extra comma
	_, err := plugins.Unmarshal[api_v1.PluginResponse](req)
	require.Error(t, err)
}

func TestMarshal_WithValidObject_ReturnsString(t *testing.T) {
	obj := api_v1.PluginResponse{
		Result: "result",
		Error:  "",
	}
	result, err := plugins.Marshal[api_v1.PluginResponse](&obj)
	expected := `{"Result":"result","Error":""}`
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestMarshal_WithInValidObject_ReturnsError(t *testing.T) {
	impossible := math.Inf(1)
	_, err := plugins.Marshal[float64](&impossible)
	assert.Equal(t, "json: unsupported value: +Inf", err.Error())
}
