package cfront

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

const (
	maxRetryCount = 120
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
		log.Println("[target path group]", pathGroup)
		callerReference := newCallerReference()
		quantity := int32(len(pathGroup))

		createInput := cloudfront.CreateInvalidationInput{
			DistributionId: &setting.DistributionId,
			InvalidationBatch: &types.InvalidationBatch{
				CallerReference: &callerReference,
				Paths: &types.Paths{
					Quantity: &quantity,
					Items:    pathGroup,
				},
			},
		}
		createOutput, err := client.CreateInvalidation(context.TODO(), &createInput)
		if err != nil {
			return err
		}
		if createOutput == nil || createOutput.Invalidation == nil || createOutput.Invalidation.Status == nil {
			return errors.New("create invalidation response status is empty")
		}

		var retryCount int
		getInput := cloudfront.GetInvalidationInput{
			DistributionId: &setting.DistributionId,
			Id:             createOutput.Invalidation.Id,
		}
	retry:
		time.Sleep(1 * time.Second)
		getOutput, err := client.GetInvalidation(context.TODO(), &getInput)
		if err != nil {
			return err
		}
		if getOutput == nil || getOutput.Invalidation == nil || getOutput.Invalidation.Status == nil {
			return errors.New("get invalidation response status is empty")
		}

		if *getOutput.Invalidation.Status != "Completed" {
			retryCount++
			if retryCount > maxRetryCount {
				return errors.New("max retry count over")
			}
			log.Println("[retry]", retryCount, "[current status]", *getOutput.Invalidation.Status)
			goto retry
		}
	}

	return nil
}

func newCallerReference() string {
	utc := time.Now().UTC()
	return utc.Format("20060102150405") // YYYYMMDDhhmmss
}
