package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
)

func MakeAwsConfig(opt Option) (cfg aws.Config, err error) {
	if opt.Profile == "" {
		cfg, err = config.LoadDefaultConfig(context.TODO())
	} else {
		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithSharedConfigProfile(opt.Profile),
			config.WithAssumeRoleCredentialOptions(func(aro *stscreds.AssumeRoleOptions) {
				aro.TokenProvider = func() (string, error) {
					return stscreds.StdinTokenProvider()
				}
			}),
		)
	}

	if err != nil {
		err = fmt.Errorf("failed to load config: %s", err)
	}
	return
}
