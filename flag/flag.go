package flag

import (
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
