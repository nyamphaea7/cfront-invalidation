package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

type CFront struct {
	Opt *Option
}

func (cf *CFront) GetOption() (*Option, error) {
	if cf.Opt == nil {
		return nil, errors.New("no set option")
	}
	return cf.Opt, nil
}

func (cf *CFront) CreateInvalidation(cfg aws.Config) error {
	opt, err := cf.GetOption()
	if err != nil {
		return err
	}

	client := cloudfront.NewFromConfig(cfg)

	for _, paths := range opt.GetPaths() {
		invalidationId, err := cf.createInvalidation(client, paths)
		if err != nil {
			return err
		}

		err = cf.getInvalidation(client, invalidationId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (cf *CFront) createInvalidation(client *cloudfront.Client, paths Paths) (*string, error) {
	opt, err := cf.GetOption()
	if err != nil {
		return nil, err
	}

	log.Println("[invalidation target paths]:", paths)
	callerReference := newCallerReference()
	quantity := int32(len(paths))

	createInput := cloudfront.CreateInvalidationInput{
		DistributionId: &opt.DistributionId,
		InvalidationBatch: &types.InvalidationBatch{
			CallerReference: &callerReference,
			Paths: &types.Paths{
				Quantity: &quantity,
				Items:    paths,
			},
		},
	}
	createOutput, err := client.CreateInvalidation(context.TODO(), &createInput)
	if err != nil {
		return nil, fmt.Errorf("failed to create invalidation: %s", err)
	}
	if createOutput == nil || createOutput.Invalidation == nil || createOutput.Invalidation.Status == nil || createOutput.Invalidation.Id == nil {
		return nil, errors.New("create invalidation response status or id is empty")
	}
	return createOutput.Invalidation.Id, nil
}

func (cf *CFront) getInvalidation(client *cloudfront.Client, invalidationId *string) error {
	opt, err := cf.GetOption()
	if err != nil {
		return err
	}

	var retryCount int
	getInput := cloudfront.GetInvalidationInput{
		DistributionId: &opt.DistributionId,
		Id:             invalidationId,
	}

retry:
	time.Sleep(time.Duration(opt.RetryInterval) * time.Second)

	getOutput, err := client.GetInvalidation(context.TODO(), &getInput)
	if err != nil {
		return err
	}
	if getOutput == nil || getOutput.Invalidation == nil || getOutput.Invalidation.Status == nil {
		return errors.New("get invalidation response status is empty")
	}

	if *getOutput.Invalidation.Status != "Completed" {
		retryCount++
		if opt.MaxRetryCount > 0 && retryCount > opt.MaxRetryCount {
			return errors.New("max retry count over")
		}

		log.Println("[retry]", retryCount, "[current status]", *getOutput.Invalidation.Status)
		goto retry
	}

	return nil
}

func newCallerReference() string {
	utc := time.Now().UTC()
	return utc.Format("20060102150405") // YYYYMMDDhhmmss
}
