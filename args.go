package main

import (
	"errors"
	"flag"
)

var (
	isSpecifiedSettingJSONFilePath bool
	isSpecifiedSomeOptions         bool

	errInvalidSpecifiedSettingFilePathAndOthers = errors.New("-s option cannot be specified with any other option(except -f). Or either -s option or other options(except -f) must be specified.")
)

func visit(f *flag.Flag) {
	if len(f.Value.String()) == 0 {
		return
	}

	switch f.Name {
	case optionNameOfSettingJSONFilePath:
		isSpecifiedSettingJSONFilePath = true
	case optionNameOfForce:
		isSpecifiedSomeOptions = true
	case optionNameOfProfile, optionNameOfDistributionId, optionNameOfPaths, optionNameOfRetryInterval, optionNameOfMaxRetryCount:
		isSpecifiedSomeOptions = true
	}
}

func GetOptions() error {
	flag.Parse()
	flag.Visit(visit)

	if isSpecifiedSettingJSONFilePath == isSpecifiedSomeOptions {
		return errInvalidSpecifiedSettingFilePathAndOthers
	}
	if isSpecifiedSettingJSONFilePath {
		if err := option.ParseOptionsFromJson(); err != nil {
			return err
		}
		return nil
	}

	return option.Valid()
}
