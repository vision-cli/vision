package config

import (
	"bufio"
	"crypto/rand"
	"errors"
	"math/big"
	"os"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/vision-cli/vision/cli"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyz")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		r, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			panic(err)
		}
		b[i] = letters[r.Int64()]
	}
	return string(b)
}

func LoadConfig(flagSet *pflag.FlagSet, silent bool, projectConfigFile, projectName string) error {
	setConfigFile(projectConfigFile)

	if !exists(projectConfigFile + configExtension) {
		cli.Warningf("Project config file %s doesnt exist", projectConfigFile)
		if silent {
			return LoadDefaultConfig(flagSet, true, projectConfigFile, projectName)
		}
		if cli.Confirmed(bufio.NewReader(os.Stdin), "Do you want to create it?") {
			return LoadDefaultConfig(flagSet, false, projectConfigFile, projectName)
		} else {
			return errors.New("project config file not found")
		}
	}

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	if viper.GetString(uniqueStr) == "" {
		viper.Set(uniqueStr, randSeq(uniqueStrLen))
	}

	return nil
}

func LoadDefaultConfig(flagSet *pflag.FlagSet, silent bool, projectConfigFile, projectName string) error {
	setConfigFile(projectConfigFile)
	if err := viper.WriteConfigAs(projectConfigFile + configExtension); err != nil {
		return err
	}

	viper.Set(uniqueStr, randSeq(uniqueStrLen))

	reader := bufio.NewReader(os.Stdin)

	genericSetter(reader, "Project name:", projectName, silent, SetProjectName)
	genericSetter(reader, "Template version:", defaultTempalateVersion, silent, SetTemplateVersion)
	mustSetWithFlag(reader, "Default remote (e.g. github.com/<company-name>/<project>):", "", silent, flagSet, FlagRemote, SetRemote)
	mustSetWithFlag(reader, "Default deployment:", defaultDeployment, silent, flagSet, FlagDeployment, SetDeployment)
	mustSetWithFlag(reader, "Default service namespace:", defaultDefaultNamespace, silent, flagSet, FlagNamespace, SetDefaultNamespace)
	mustSetWithFlag(reader, "Default registry:", defaultRegistry, silent, flagSet, FlagRegistry, SetRegistry)
	mustSetWithFlag(reader, "Default branch:", defaultBranch, silent, flagSet, FlagBranch, SetBranch)
	mustSetWithFlag(reader, "Default api version:", defaultApiVersion, silent, flagSet, FlagApiVersion, SetApiVersion)
	genericSetter(reader, "Default gateway service name:", defaultGatewayName, silent, SetGatewayName)
	genericSetter(reader, "Default graphql service name:", defaultGraphqlName, silent, SetGraphqlName)
	genericSetter(reader, "Default services directory:", defaultServicesDir, silent, SetServicesDirectory)
	genericSetter(reader, "Default infra directory:", defaultInfraDir, silent, SetInfraDirectory)

	return nil
}

func genericSetter(reader *bufio.Reader, prompt, def string, slient bool, setter func(string)) {
	if slient {
		setter(def)
	} else {
		ans := cli.Input(reader, prompt, def, true)
		setter(ans)
	}
}

func mustSetWithFlag(
	reader *bufio.Reader,
	prompt, def string,
	slient bool,
	flagSet *pflag.FlagSet,
	flagName string,
	setter func(string)) {
	if strval, err := flagSet.GetString(flagName); err == nil && strval != "" {
		setter(strval)
		return
	}
	val := def
	if !slient {
		val = cli.Input(reader, prompt, val, true)
	}
	if val == "" {
		cli.Fatalf("Value of %s cannot be empty", flagName)
	}
	setter(val)
}

func SaveConfig() error {
	return viper.WriteConfig()
}

func setConfigFile(projectConfigFile string) {
	viper.SetConfigName(projectConfigFile)
	viper.AddConfigPath(".")
}

func setAndSave(key, value string) {
	viper.Set(key, value)
	if err := SaveConfig(); err != nil {
		panic(err)
	}
}

func exists(path string) bool {
	_, err := os.Stat(path)
	// other errors or nil imply existence (e.g. ErrPermission)
	return !(errors.Is(err, os.ErrNotExist) || errors.Is(err, os.ErrInvalid))
}
