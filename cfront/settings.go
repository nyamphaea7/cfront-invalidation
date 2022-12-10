package cfront

import (
	"context"
	"encoding/json"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type Setting struct {
	InvalidationGroupId string     `json:"invalidationGroupId"`
	Profile             string     `json:"profile"`
	DistributionId      string     `json:"distributionId"`
	PathGroups          [][]string `json:"pathGroups"`
}

func (s *Setting) GetAwsConfig() (aws.Config, error) {
	if s.Profile == "" {
		return config.LoadDefaultConfig(context.TODO())
	}
	return config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(s.Profile),
	)
}

type Settings []Setting

const defaultSettingsJsonFilePath = "./settings.json"

var settings Settings

func init() {
	settingsJsonFilePath := os.Getenv("SETTINGS_JSON_FILE_PATH")
	if settingsJsonFilePath == "" {
		settingsJsonFilePath = defaultSettingsJsonFilePath
	}
	_, err := os.Stat(settingsJsonFilePath)
	if err != nil {
		panic(err)
	}

	raw, err := os.ReadFile(settingsJsonFilePath)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(raw, &settings)
	if err != nil {
		panic(err)
	}
}

func GetSettings(ids InvalidationGroupIds) Settings {
	if len(ids) == 0 {
		return settings
	}

	var filtered Settings
	for _, setting := range settings {
		if contain(ids, setting.InvalidationGroupId) {
			filtered = append(filtered, setting)
		}
	}
	return filtered
}
