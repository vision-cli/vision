package config

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/viper"
)

const (
	viperAzurePrefix            = "azure"
	viperAzureResourceGroupKey  = viperAzurePrefix + ".resource-group"
	viperAzureLocationKey       = viperAzurePrefix + ".location"
	viperAzureStorageAccountKey = viperAzurePrefix + ".storage-account"
	viperAzureTennantKey        = viperAzurePrefix + ".tenant"
	viperAzureKeyvaultKey       = viperAzurePrefix + ".keyvault"
	viperAzureAcrNameKey        = viperAzurePrefix + ".acr.name"
	viperAzureAcrLoginServerKey = viperAzurePrefix + ".acr.login-server"
	maxAzureKeyLen              = 20
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9 ]+`)

func cleanString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func namify(key string, maxLen int) string { //nolint:unparam // maxLen is used during testing
	name := fmt.Sprintf("%s-%s-%s", cleanString(ProjectName()), key, viper.GetString(uniqueStr))
	if len(name) > namifyLen {
		name = name[:maxLen]
	}
	return name
}

func plainNamify(key string) string {
	return strings.ReplaceAll(namify(key, maxAzureKeyLen), "-", "")
}

func DefaultAzureResourceGroupName() string {
	return namify("rg", maxAzureKeyLen)
}

func AzureResourceGroup() string {
	return viper.GetString(viperAzureResourceGroupKey)
}

func SetAzureResourceGroup(s string) {
	setAndSave(viperAzureResourceGroupKey, s)
}

func DefaultAzureLocation() string {
	return "uksouth"
}

func AzureLocation() string {
	return viper.GetString(viperAzureLocationKey)
}

func SetAzureLocation(s string) {
	setAndSave(viperAzureLocationKey, s)
}

func DefaultAzureStorageAccount() string {
	return plainNamify("sa")
}

func AzureStorageAccount() string {
	return viper.GetString(viperAzureStorageAccountKey)
}

func SetAzureStorageAccount(s string) {
	setAndSave(viperAzureStorageAccountKey, s)
}

func AzureTenant() string {
	return viper.GetString(viperAzureTennantKey)
}

func SetAzureTenant(s string) {
	setAndSave(viperAzureTennantKey, s)
}

func DefaultAzureKeyvault() string {
	return namify("kv", maxAzureKeyLen)
}

func AzureKeyvault() string {
	return viper.GetString(viperAzureKeyvaultKey)
}

func SetAzureKeyvault(s string) {
	setAndSave(viperAzureKeyvaultKey, s)
}

func DefaultAzureAcr() string {
	return plainNamify("cr")
}

func AzureAcr() string {
	return viper.GetString(viperAzureAcrNameKey)
}

func SetAzureAcr(s string) {
	setAndSave(viperAzureAcrNameKey, s)
}

func AzureAcrLoginServer() string {
	return viper.GetString(viperAzureAcrLoginServerKey)
}

func SetAzureAcrLoginServer(s string) {
	setAndSave(viperAzureAcrLoginServerKey, s)
}

func DefaultAzureApp() string {
	return ProjectName() + "-" + viper.GetString(uniqueStr)
}
