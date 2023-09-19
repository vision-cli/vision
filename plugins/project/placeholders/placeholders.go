package placeholders

import (
	"net/url"
	"regexp"

	"github.com/barkimedes/go-deepcopy"
	api_v1 "github.com/vision-cli/api/v1"
)

const (
	ArgsCommandIndex = 0
	ArgsNameIndex    = 1
	// include any other arg indexes here
)

var allowedRegex = regexp.MustCompile(`[^a-zA-Z\-\_]+`)

func SetupPlaceholders(req api_v1.PluginRequest) (*api_v1.PluginPlaceholders, error) {
	var err error
	p, err := deepcopy.Anything(&req.Placeholders)
	if err != nil {
		return nil, err
	}
	projectName := clearName(req.Args[ArgsNameIndex])
	p.(*api_v1.PluginPlaceholders).ProjectRoot = projectName
	p.(*api_v1.PluginPlaceholders).ProjectName = projectName
	p.(*api_v1.PluginPlaceholders).ProjectDirectory = projectName
	p.(*api_v1.PluginPlaceholders).ProjectFqn, err = url.JoinPath(req.Placeholders.Remote, projectName)
	if err != nil {
		return nil, err
	}
	p.(*api_v1.PluginPlaceholders).LibsFqn, err = url.JoinPath(req.Placeholders.Remote, projectName, "libs")
	if err != nil {
		return nil, err
	}
	return p.(*api_v1.PluginPlaceholders), nil
}

func clearName(str string) string {
	return allowedRegex.ReplaceAllString(str, "")
}
