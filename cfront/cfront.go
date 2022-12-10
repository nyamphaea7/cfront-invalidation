package cfront

import (
	"context"
	"errors"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

func CreateInvalidations(ids InvalidationGroupIds) error {
	settings := GetSettings(ids)
	if len(settings) == 0 {
		return errors.New("not found InvalidationGroupIds in settings")
	}

	for _, setting := range settings {
		err := createInvalidation(setting)
		if err != nil {
			return err
		}
	}

	return nil
}

func createInvalidation(setting Setting) error {
	cfg, err := setting.GetAwsConfig()
	if err != nil {
		return err
	}
	client := cloudfront.NewFromConfig(cfg)

	for _, pathGroup := range setting.PathGroups {
		callerReference := newCallerReference()
		quantity := int32(len(pathGroup))

		input := cloudfront.CreateInvalidationInput{
			DistributionId: &setting.DistributionId,
			InvalidationBatch: &types.InvalidationBatch{
				CallerReference: &callerReference,
				Paths: &types.Paths{
					Quantity: &quantity,
					Items:    pathGroup,
				},
			},
		}
	retry:
		output, err := client.CreateInvalidation(context.TODO(), &input)
		if err != nil {
			return err
		}
		if output == nil || output.Invalidation == nil || output.Invalidation.Status == nil {
			return errors.New("create invalidation response status is empty")
		}
		if *output.Invalidation.Status != "Completed" {
			time.Sleep(1 * time.Second)
			goto retry
		}
	}

	return nil
}

func newCallerReference() string {
	utc := time.Now().UTC()
	return utc.Format("20060102150405") // YYYYMMDDhhmmss
}
