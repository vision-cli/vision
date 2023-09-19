package marshal_test

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api_v1 "github.com/vision-cli/api/v1"
	"github.com/vision-cli/vision/common/marshal"
)

func TestUnmarshal_WithValidInputProvided_ReturnsObject(t *testing.T) {
	req := `{"Result":"result","Error":""}`
	result, err := marshal.Unmarshal[api_v1.PluginResponse](req)
	expected := api_v1.PluginResponse{
		Result: "result",
		Error:  "",
	}
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestUnmarshal_WhenInValidInputProvided_ReturnsError(t *testing.T) {
	req := `{"Result":"result","Error":"",}` // extra comma
	_, err := marshal.Unmarshal[api_v1.PluginResponse](req)
	require.Error(t, err)
}

func TestUnmarshal_WhenFieldMissing_ReturnsError(t *testing.T) {
	req := `{"Something":"else"}`
	_, err := marshal.Unmarshal[api_v1.PluginResponse](req)
	require.Error(t, err)
	assert.Equal(t, `json: unknown field "Something"`, err.Error())
}

func TestMarshal_WithValidObject_ReturnsString(t *testing.T) {
	obj := api_v1.PluginResponse{
		Result: "result",
		Error:  "",
	}
	result, err := marshal.Marshal[api_v1.PluginResponse](obj)
	expected := `{"Result":"result","Error":""}`
	require.NoError(t, err)
	assert.Equal(t, expected, result)
}

func TestMarshal_WithInValidObject_ReturnsError(t *testing.T) {
	_, err := marshal.Marshal[float64](math.Inf(1))
	assert.Equal(t, "json: unsupported value: +Inf", err.Error())
}
