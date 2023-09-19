package plugin_test

import (
	"io/fs"
	"strings"
	"testing"

	"github.com/vision-cli/vision/plugins/project/plugin"
	"github.com/vision-cli/vision/plugins/project/run"

	"github.com/stretchr/testify/assert"
	"github.com/vision-cli/vision/common/file"
	"github.com/vision-cli/vision/common/mocks"
)

func TestHandle_WithValidUsageInput_ReturnsUsageResponseString(t *testing.T) {
	e := mocks.NewMockExecutor()
	tw := mocks.NewMockTmplWriter()
	result := plugin.Handle(CreateRequest(t, "usage"), &e, &tw)
	expected := `{"Version":"0.1.0","Use":"project","Short":"manage projects","Long":"manage projects and docs using a standard template","Example":"vision project create myProject  -r github.com/mycompany","Subcommands":["create"],"Flags":[],"RequiresConfig":false}`
	assert.Equal(t, expected, result)
}

func TestHandle_WithValidConfigInput_ReturnsConfigResponseString(t *testing.T) {
	e := mocks.NewMockExecutor()
	tw := mocks.NewMockTmplWriter()
	result := plugin.Handle(CreateRequest(t, "config"), &e, &tw)
	expected := `{"Defaults":[]}`
	assert.Equal(t, expected, result)
}

func TestHandle_WithInValidInput_ReturnsErrorString(t *testing.T) {
	e := mocks.NewMockExecutor()
	tw := mocks.NewMockTmplWriter()
	result := plugin.Handle("X"+CreateRequest(t, "run"), &e, &tw)
	expected := `{"Result":"","Error":"invalid character 'X' looking for beginning of value"}`
	assert.Equal(t, expected, result)
}

func TestHandle_WithInValidCommand_ReturnsErrorString(t *testing.T) {
	e := mocks.NewMockExecutor()
	tw := mocks.NewMockTmplWriter()
	req := CreateRequest(t, "avengers")
	result := plugin.Handle(req, &e, &tw)
	expected := `{"Result":"","Error":"unknown command"}`
	assert.Equal(t, expected, result)
}

func TestHandle_WithValidRunInput_ReturnsRunResponseString(t *testing.T) {
	oldtolines := file.ToLines
	defer func() { file.ToLines = oldtolines }()
	file.ToLines = func(path string) ([]string, error) {
		return []string{"# Docs"}, nil
	}
	oldfromlines := file.FromLines
	defer func() { file.FromLines = oldfromlines }()
	file.FromLines = func(path string, lines []string) error {
		return nil
	}
	oldfswalkdir := run.Fswalkdir
	defer func() { run.Fswalkdir = oldfswalkdir }()
	run.Fswalkdir = func(fsys fs.FS, root string, fn fs.WalkDirFunc) error {
		return nil
	}

	e := mocks.NewMockExecutor()
	tw := mocks.NewMockTmplWriter()
	req := CreateRequest(t, "run")
	req = strings.Replace(req, `"Args":[]`, `"Args":["create","mything"]`, 1)
	result := plugin.Handle(req, &e, &tw)
	expected := `{"Result":"SUCCESS!","Error":""}`
	assert.Equal(t, expected, result)
}

func CreateRequest(t *testing.T, command string) string {
	t.Helper()
	var testReq = `
{
	"Command":"` + command + `",
	"Args":[],
	"Flags":[],
	"Placeholders": {
		"ProjectRoot":"",
		"ProjectName":"",
		"ProjectDirectory":"",
		"ProjectFqn":"",
		"Registry":"",
		"Remote":"",
		"Branch":"",
		"Version":"",
		"ServicesFqn":"",
		"ServicesDirectory":"",
		"GatewayServiceName":"",
		"GatewayFqn":"",
		"GraphqlServiceName":"",
		"GraphqlFqn":"",
		"LibsFqn":"",
		"LibsDirectory":"",
		"ServiceNamespace":"",
		"ServiceVersionedNamespace":"",
		"ServiceName":"",
		"ServiceFqn":"",
		"ServiceDirectory":"",
		"InfraDirectory":"",
		"ProtoPackage":""
		}
}	
`
	return testReq
}
