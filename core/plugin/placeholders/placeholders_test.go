package placeholders_test

import (
	"testing"

	"github.com/vision-cli/vision/core/plugin/placeholders"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api_v1 "github.com/vision-cli/api/v1"
)

func TestSetupPlaceholders_WithValidNameAndVersion_ReturnsCorrectPlaceholders(t *testing.T) {
	r := api_v1.PluginRequest{
		Args: []string{"create", "myplugin", "v2"},
		Placeholders: api_v1.PluginPlaceholders{
			Remote: "github.com/mycompany",
		},
	}
	result, err := placeholders.SetupPlaceholders(r)
	require.NoError(t, err)
	expected := &placeholders.Placeholders{
		Name:      "myplugin",
		Namespace: "github.com/mycompany/vision-plugin-myplugin-v2",
		Directory: "vision-plugin-myplugin-v2",
	}
	assert.Equal(t, expected, result)
}

func TestSetupPlaceholders_WithInvalidName_ReturnsError(t *testing.T) {
	r := api_v1.PluginRequest{
		Args: []string{"create", "plugin", "v2"},
	}
	_, err := placeholders.SetupPlaceholders(r)
	require.Error(t, err)
}

func TestSetupPlaceholders_WithVersionMissingV_ReturnsError(t *testing.T) {
	r := api_v1.PluginRequest{
		Args: []string{"create", "plugin", "2"},
	}
	_, err := placeholders.SetupPlaceholders(r)
	require.Error(t, err)
}

func TestSetupPlaceholders_WithInvalidVersion_ReturnsError(t *testing.T) {
	r := api_v1.PluginRequest{
		Args: []string{"create", "plugin", "v2.0"},
	}
	_, err := placeholders.SetupPlaceholders(r)
	require.Error(t, err)
}
