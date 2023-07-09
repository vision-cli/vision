package plugins_test

import (
	"os"
	"os/exec"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vision-cli/common/mocks"
	commonPlugins "github.com/vision-cli/common/plugins"
	"github.com/vision-cli/vision/config"
	"github.com/vision-cli/vision/plugins"
)

var usageResp = `{
	"Version":        "0.1.0",
	"Use":            "project",
	"Short":          "project short description",
	"Long":           "project long description",
	"Example":        "project example",
	"Subcommands":    ["create", "delete"],
	"Flags":          [],
	"RequiresConfig": false
}`

var projectexe = commonPlugins.Plugin{
	Name:            "projectexe",
	PluginPath:      "projectexe",
	InternalCommand: nil,
}

func TestGetCobraCommand_WithValidUsageResponse_ReturnsCobraCommand(t *testing.T) {
	te := mocks.NewMockExecutor()
	cmd := GetGoodCmd(t, &te)
	assert.Equal(t, "project", cmd.Use)
	assert.Equal(t, "project", cmd.Use)
	assert.Equal(t, "project short description", cmd.Short)
	assert.Equal(t, 1, len(te.History()))
}

func TestGetCobraCommand_WithInvalidUsageResponse_ReturnsError(t *testing.T) {
	te := mocks.NewMockExecutor()
	te.SetOutput("X" + usageResp)
	_, err := plugins.GetCobraCommand(projectexe, &te)
	require.Error(t, err)
}

func TestReturnedCobraCommand_WithoutRequiredConfigAndRunSuccess_ReturnsSuccess(t *testing.T) {
	te := mocks.NewMockExecutor()
	cmd := GetGoodCmd(t, &te)

	te.SetOutput(`{"Result":"Success","Error":""}`)
	res := mocks.WithMockStdout(t, func() {
		cmd.Run(cmd, []string{})
	})

	assert.Equal(t, "Success\n", res[len(res)-8:])
	assert.Equal(t, 2, len(te.History()))
}

func TestReturnedCobraCommand_WithRequiredConfigAndRunSuccess_OsExitNoConfigFile(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		te := mocks.NewMockExecutor()
		var usageRespWithConfig = `{
			"Version":        "0.1.0",
			"Use":            "project",
			"Short":          "project short description",
			"Long":           "project long description",
			"Example":        "project example",
			"Subcommands":    ["create", "delete"],
			"Flags":          [],
			"RequiresConfig": true
		}`

		te.SetOutput(usageRespWithConfig)
		cmd, err := plugins.GetCobraCommand(projectexe, &te)
		require.NoError(t, err)

		cmd.Run(cmd, []string{})
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestReturnedCobraCommand_WithRequiredConfigAndRunSuccess_OsExitNoConfigFile")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		assert.Equal(t, 1, e.ExitCode())
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)

}

func TestReturnedCobraCommand_WithoutRequiredFlag_OsExit(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		te := mocks.NewMockExecutor()
		cmd := GetBadCmd(t, &te)
		cmd.Run(cmd, []string{})
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestReturnedCobraCommand_WithoutRequiredFlag_OsExit")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		assert.Equal(t, 1, e.ExitCode())
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestReturnedCobraCommand_WithoutRequiredConfigAndCallFailure_OsExit(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		te := mocks.NewMockExecutor()
		cmd := GetGoodCmd(t, &te)
		te.SetOutput(`X{"Result":"Success","Error":""}`)
		cmd.Run(cmd, []string{})
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestReturnedCobraCommand_WithoutRequiredConfigAndCallFailure_OsExit")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		assert.Equal(t, 1, e.ExitCode())
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func TestReturnedCobraCommand_WithoutRequiredConfigAndCallSuccessButMessageError_OsExit(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		te := mocks.NewMockExecutor()
		cmd := GetGoodCmd(t, &te)
		te.SetOutput(`{"Result":"","Error":"Error message"}`)
		cmd.Run(cmd, []string{})
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestReturnedCobraCommand_WithoutRequiredConfigAndCallSuccessButMessageError_OsExit")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		assert.Equal(t, 1, e.ExitCode())
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

func GetGoodCmd(t *testing.T, te *mocks.MockExecutor) *cobra.Command {
	t.Helper()
	cmd := GetBadCmd(t, te)
	err := cmd.Flags().Lookup(config.FlagRemote).Value.Set("github.com/mycompany")
	require.NoError(t, err)
	cmd.Flags().Lookup(config.FlagRemote).Changed = true
	return cmd
}

func GetBadCmd(t *testing.T, te *mocks.MockExecutor) *cobra.Command {
	t.Helper()
	te.SetOutput(usageResp)
	cmd, err := plugins.GetCobraCommand(projectexe, te)
	require.NoError(t, err)
	return cmd
}
