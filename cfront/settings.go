package cfront

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
)

type Setting struct {
	InvalidationGroupId string     `json:"invalidationGroupId"`
	Profile             string     `json:"profile"`
	DistributionId      string     `json:"distributionId"`
	PathGroups          [][]string `json:"pathGroups"`
}

func (s *Setting) IsValid() bool {
	if s.DistributionId == "" {
		return false
	}

	if len(s.PathGroups) == 0 {
		return false
	}

	var isValidSomePathGroups bool
	for _, group := range s.PathGroups {
		isValidSomePathGroups = len(group) > 0
		if isValidSomePathGroups {
			break
		}
	}

	return isValidSomePathGroups
}

func (s *Setting) GetAwsConfig() (aws.Config, error) {
	if s.Profile == "" {
		return config.LoadDefaultConfig(context.TODO())
	}
	return config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(s.Profile),
		config.WithAssumeRoleCredentialOptions(func(aro *stscreds.AssumeRoleOptions) {
			aro.TokenProvider = func() (string, error) {
				return stscreds.StdinTokenProvider()
			}
		}),
	)
}

type Settings []Setting

const defaultSettingsJsonFilePath = "./settings.json"

var baseSettings Settings

func init() {
	settingsJsonFilePath := os.Getenv("SETTINGS_JSON_FILE_PATH")
	if settingsJsonFilePath == "" {
		settingsJsonFilePath = defaultSettingsJsonFilePath
	}
	_, err := os.Stat(settingsJsonFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	raw, err := os.ReadFile(settingsJsonFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	err = json.Unmarshal(raw, &baseSettings)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}

	var filtered Settings
	for _, setting := range baseSettings {
		if setting.IsValid() {
			filtered = append(filtered, setting)
		}
	}
	baseSettings = filtered

	if len(baseSettings) == 0 {
		fmt.Println("no settings are valid. please confirm settings.")
		os.Exit(2)
	}
}

func GetSettings(ids InvalidationGroupIds) Settings {
	if len(ids) == 0 {
		return baseSettings
	}

	var filtered Settings
	for _, setting := range baseSettings {
		if contain(ids, setting.InvalidationGroupId) {
			filtered = append(filtered, setting)
		}
	}
	return filtered
}
