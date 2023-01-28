package main

import (
	"flag"
	"fmt"
	"os"
)

const (
	optionNameOfForce = "f"

	optionNameOfProfile        = "p"
	optionNameOfDistributionId = "d"
	optionNameOfPaths          = "t"
	optionNameOfRetryInterval  = "i"
	optionNameOfMaxRetryCount  = "c"

	optionNameOfSettingJSONFilePath = "s"

	defaultRetryInterval = 1
	defaultMaxRetryCount = 60
)

var option Option

func init() {
	flag.BoolVar(
		&option.Force,
		optionNameOfForce,
		false,
		"Start creating invalidation no confirmation of profile, distribution id and paths.",
	)

	flag.StringVar(
		&option.Profile,
		optionNameOfProfile,
		"",
		"AWS profile name to use to create invalidation. Automatically, the default profile use unless -p is specified.",
	)
	flag.StringVar(
		&option.DistributionId,
		optionNameOfDistributionId,
		"",
		"Invalidation target Amazon CloudFront distribution id. Automatically, the default profile use unless -p is specified.",
	)
	flag.Var(
		&option.Paths,
		optionNameOfPaths,
		"Invalidation target paths.\nspecified example:\n\t-t /index.html -t /foo/bar/*\n\t-t /index.html,/foo/bar/*",
	)

	flag.IntVar(
		&option.RetryInterval,
		optionNameOfRetryInterval,
		defaultRetryInterval,
		"Interval time(sec) to check status of creating invalidation. Must be at least 1 second.",
	)
	flag.IntVar(
		&option.MaxRetryCount,
		optionNameOfMaxRetryCount,
		defaultMaxRetryCount,
		"Max count to retry check status of creating invalidation. Specify 0 to ignore the max count of retries and execute infinitely.",
	)

	flag.StringVar(
		&option.SettingJSONFilePath,
		optionNameOfSettingJSONFilePath,
		"",
		"Settings JSON file path. This option cannot be specified with any other options(except -f).",
	)

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "\nUsage: %s [options...]\n\nOptions:\n", os.Args[0])
		flag.PrintDefaults()
	}
}
