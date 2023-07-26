package config

import (
	"bufio"
	"errors"
	"fmt"
	"os"

	"github.com/spf13/pflag"

	"github.com/vision-cli/vision/cli"
)

var stat = os.Stat
var v = Persist(&ViperPersist{})

func LoadConfig(flagSet *pflag.FlagSet, silent bool, projectConfigFile, projectName string, requireConfig bool) error {
	setConfigFile(projectConfigFile)
	reader := bufio.NewReader(os.Stdin)

	if !exists(projectConfigFile+configExtension) && requireConfig {
		cli.Warningf("Project config file %s doesnt exist", projectConfigFile)
		if silent {
			return loadDefaultConfig(flagSet, true, projectConfigFile, projectName, reader)
		}
		if cli.Confirmed(reader, "Do you want to create it?") {
			return loadDefaultConfig(flagSet, false, projectConfigFile, projectName, reader)
		}
		return fmt.Errorf("project config file not found")
	}

	err := v.ReadInConfig()
	if err != nil {
		return err
	}

	if v.GetString(uniqueStr) == "" {
		v.Set(uniqueStr, randSeq(uniqueStrLen))
	}

	return nil
}

func loadDefaultConfig(
	flagSet *pflag.FlagSet,
	silent bool,
	projectConfigFile, projectName string,
	reader *bufio.Reader) error {
	if err := v.WriteConfigAs(projectConfigFile + configExtension); err != nil {
		return err
	}

	v.Set(uniqueStr, randSeq(uniqueStrLen))
	genericSetter(reader, "Project name:", projectName, silent, SetProjectName)

	for _, def := range defaultConfigsWithFlags {
		err := mustSetWithFlag(reader, def.prompt, def.def, silent, flagSet, def.flagname, def.setter)
		if err != nil {
			return err
		}
	}

	for _, def := range defaultConfigs {
		genericSetter(reader, def.prompt, def.def, silent, def.setter)
	}

	return nil
}

func genericSetter(reader *bufio.Reader, prompt, def string, silent bool, setter func(string)) {
	if silent {
		setter(def)
	} else {
		ans := cli.Input(reader, prompt, def, true)
		setter(ans)
	}
}

func mustSetWithFlag(
	reader *bufio.Reader,
	prompt, def string,
	silent bool,
	flagSet *pflag.FlagSet,
	flagName string,
	setter func(string)) error {
	if strval, err := flagSet.GetString(flagName); err == nil && strval != "" {
		setter(strval)
		return nil
	}
	val := def
	if !silent {
		val = cli.Input(reader, prompt, val, true)
	}
	if val == "" {
		return fmt.Errorf("value of %s cannot be empty", flagName)
	}
	setter(val)
	return nil
}

func setConfigFile(projectConfigFile string) {
	v.SetConfigName(projectConfigFile)
	v.AddConfigPath(".")
}

func setAndSave(key, value string) {
	v.Set(key, value)
	if err := v.WriteConfig(); err != nil {
		panic(err)
	}
}

func exists(path string) bool {
	_, err := stat(path)
	// other errors or nil imply existence (e.g. ErrPermission)
	return !(errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrInvalid))
}

func UniqueStr() string {
	return v.GetString(uniqueStr)
}
