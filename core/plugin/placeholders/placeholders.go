package placeholders

import (
	"fmt"
	"net/url"
	"regexp"

	api_v1 "github.com/vision-cli/api/v1"
)

const (
	ArgsCommandIndex = 0
	ArgsNameIndex    = 1
	ArgsVersionIndex = 2
)

var nonAlphaRegex = regexp.MustCompile(`[^a-zA-Z]+`)
var validVersion = regexp.MustCompile(`^v\d+$`)

type Placeholders struct {
	Name      string
	Namespace string
	Directory string
}

func SetupPlaceholders(req api_v1.PluginRequest) (*Placeholders, error) {
	name := clearString(req.Args[ArgsNameIndex])
	version := req.Args[ArgsVersionIndex]

	banned := map[string]struct{}{
		"imgaes":       {},
		"placeholders": {},
		"plugin":       {},
	}

	if _, match := banned[name]; match {
		return nil, fmt.Errorf("name not alowed")
	}

	if !validVersion.MatchString(version) {
		return nil, fmt.Errorf("invalid version")
	}

	namespace, err := url.JoinPath(req.Placeholders.Remote, fmt.Sprintf("vision-plugin-%s-%s", name, version))
	if err != nil {
		return nil, err
	}

	return &Placeholders{
		Name:      name,
		Namespace: namespace,
		Directory: fmt.Sprintf("vision-plugin-%s-%s", name, version),
	}, nil
}

func clearString(str string) string {
	return nonAlphaRegex.ReplaceAllString(str, "")
}
