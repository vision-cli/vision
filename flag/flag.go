package flag

import (
	"fmt"
	"log"

	"github.com/spf13/pflag"

	"github.com/vision-cli/vision/config"
)

// MustGetString returns the value for the named string flag or exits.
func MustGetString(flagSet *pflag.FlagSet, name string) string {
	value, err := flagSet.GetString(name)
	if err != nil {
		log.Fatalf("error retrieving value for flag [%s]: %v\n", name, err)
	}
	return value
}

func IsSilent(flagSet *pflag.FlagSet) bool {
	silent, err := flagSet.GetBool(config.FlagSilent)
	if err != nil {
		log.Fatalf("error retrieving value for flag [%s]: %v\n", config.FlagSilent, err)
	}
	return silent
}

func IsForce(flagSet *pflag.FlagSet) bool {
	force, err := flagSet.GetBool(config.FlagForce)
	if err != nil {
		log.Fatalf("error retrieving value for flag [%s]: %v\n", config.FlagForce, err)
	}
	return force
}

func IsFromTemplate(flagSet *pflag.FlagSet) bool {
	template := MustGetString(flagSet, config.FlagTemplate)
	return template != ""
}

// Config returns the manadatory overridable flagset from the project's config file
func ConfigFlagset() *pflag.FlagSet {
	conf := pflag.NewFlagSet("config", 1)
	conf.StringP(config.FlagRemote, "r", config.Remote(), "remote url for project repo")
	conf.StringP(config.FlagBranch, "b", config.Branch(), "default branch in the project repo")
	conf.StringP(config.FlagRegistry, "g", config.Registry(), "container registry")
	conf.StringP(config.FlagNamespace, "n", config.DefaultNamespace(), "default namespace to use")
	conf.StringP(config.FlagApiVersion, "v", config.ApiVersion(), "api version to use")
	conf.StringP(config.FlagTemplate, "t", "", "template to use")
	conf.StringP(config.FlagDeployment, "d", config.Deployment(),
		fmt.Sprintf("deployment pattern to use [%s, %s, %s]", config.DeployPlatform, config.DeployStandaloneGraphql))
	conf.Bool(config.FlagSilent, false, "use default values for all flags and dont ask questions")
	conf.Bool(config.FlagForce, false, "overwrite without asking questions")
	return conf
}
