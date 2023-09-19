package placeholders

import (
	"regexp"

	api_v1 "github.com/vision-cli/api/v1"
)

const (
	ArgsCommandIndex = 0
	ArgsNameIndex    = 1
	// include any other arg indexes here
)

var nonAlphaRegex = regexp.MustCompile(`[^a-zA-Z]+`)

type Placeholders struct {
	Name string
}

func SetupPlaceholders(req api_v1.PluginRequest) (*Placeholders, error) {
	// setup your placeholders here
	// you can also deepcopy the Placeholders in the plugin request and use it
	// this is just an example:
	name := clearString(req.Args[ArgsNameIndex])
	return &Placeholders{
		Name: name,
	}, nil
}

func clearString(str string) string {
	return nonAlphaRegex.ReplaceAllString(str, "")
}
