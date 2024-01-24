package plugin

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// plugin executable follows the format:
// vision-plugin-[name]-[version]
//
//	[0]    [1]    [2]    [3]
const PREFIX = "vision-plugin"
const NAME_INDEX = 2
const VERSION_INDEX = 3

type Plugin struct {
	Name     string
	Version  string
	FullPath string
}

// Find searches all dirs in the PATH envar to find binaries with specific vision formatting and assigns them to a map.
// The formatting is `vision-plugin-[plugin-name]-[version-number]`
func Find() ([]Plugin, error) {
	var plugins []Plugin

	sysPath := os.Getenv("PATH")
	paths := strings.Split(sysPath, string(os.PathListSeparator))

	m := make(map[string]struct{})
	for _, p := range paths {
		m[p] = struct{}{}
	}

	for path := range m {
		err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
			if info == nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}
			if strings.HasPrefix(info.Name(), PREFIX) {
				pluginSplit := strings.Split(info.Name(), "-")
				if len(pluginSplit) == 4 { // only process correctly named plugins
					name, version := pluginSplit[2], pluginSplit[3]
					plugins = append(plugins, Plugin{
						Name:     name,
						Version:  version,
						FullPath: path,
					})
				}
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}
	return plugins, nil
}
