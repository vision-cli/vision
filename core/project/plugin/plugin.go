package plugin

import (
	"errors"

	"github.com/vision-cli/vision/core/project/placeholders"
	"github.com/vision-cli/vision/core/project/run"

	api_v1 "github.com/vision-cli/api/v1"
	"github.com/vision-cli/vision/common/execute"
	"github.com/vision-cli/vision/common/marshal"
	"github.com/vision-cli/vision/common/tmpl"
)

const (
	create = "create"
)

var Usage = api_v1.PluginUsageResponse{
	Version:        "0.1.0",
	Use:            "project",
	Short:          "manage projects",
	Long:           "manage projects and docs using a standard template",
	Example:        "vision project create myProject  -r github.com/mycompany",
	Subcommands:    []string{create},
	Flags:          []api_v1.PluginFlag{},
	RequiresConfig: false,
}

var DefaultConfig = api_v1.PluginConfigResponse{
	Defaults: []api_v1.PluginConfigItem{},
}

func Handle(input string, e execute.Executor, t tmpl.TmplWriter) string {
	req, err := marshal.Unmarshal[api_v1.PluginRequest](input)
	if err != nil {
		return errorResponse(err)
	}
	result := ""
	switch req.Command {
	case api_v1.CommandUsage:
		result, err = marshal.Marshal[api_v1.PluginUsageResponse](Usage)
		if err != nil {
			return errorResponse(err)
		}
	case api_v1.CommandConfig:
		result, err = marshal.Marshal[api_v1.PluginConfigResponse](DefaultConfig)
		if err != nil {
			return errorResponse(err)
		}
	case api_v1.CommandRun:
		if len(req.Args) == 0 || req.Args[placeholders.ArgsCommandIndex] == "" {
			return errorResponse(errors.New("missing cli command"))
		}
		switch req.Args[placeholders.ArgsCommandIndex] {
		case create:
			if len(req.Args) <= 1 ||
				req.Args[placeholders.ArgsNameIndex] == "" {
				return errorResponse(errors.New("missing project name"))
			}
			p, err := placeholders.SetupPlaceholders(req)
			if err != nil {
				return errorResponse(err)
			}
			err = run.Create(p, e, t)
			if err != nil {
				return errorResponse(err)
			}
		default:
			return errorResponse(errors.New("unknown cli command"))
		}
		resp := api_v1.PluginResponse{
			Result: "SUCCESS!",
			Error:  "",
		}
		result, err = marshal.Marshal[api_v1.PluginResponse](resp)
		if err != nil {
			return errorResponse(err)
		}
	default:
		return errorResponse(errors.New("unknown command"))
	}
	return result
}

func errorResponse(err error) string {
	res, err := marshal.Marshal[api_v1.PluginResponse](api_v1.PluginResponse{
		Result: "",
		Error:  err.Error(),
	})
	if err != nil {
		panic(err.Error())
	}
	return res
}
