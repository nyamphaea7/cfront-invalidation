package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Paths []string

func (p *Paths) String() string {
	return ""
}

func (p *Paths) Set(v string) error {
	paths := strings.Split(v, ",")
	for _, path := range paths {
		if len(path) > 1 {
			*p = append(*p, path)
		}
	}

	return nil
}

func (p *Paths) ValidateURIs() error {
	var invalidPaths []string

	for _, path := range *p {
		if _, err := url.ParseRequestURI(path); err != nil {
			invalidPaths = append(invalidPaths, path)
		}
	}

	if len(invalidPaths) > 0 {
		return fmt.Errorf(errTemplateInvalidPathURIs, invalidPaths)
	}
	return nil
}

type Option struct {
	Force bool `json:"force"`

	Profile        string `json:"profile"`
	DistributionId string `json:"distributionId"`
	Paths          Paths
	PathGroups     []Paths `json:"pathGroups"`

	RetryInterval int `json:"retryInterval"`
	MaxRetryCount int `json:"maxRetryCount"`

	SettingJSONFilePath string
}

var (
	errInvalidDistributionId = errors.New("no specified distribution id")
	errInvalidPathLength     = errors.New("no specified any paths")
	errInvalidRetryInterval  = errors.New("the retry interval time must be 1 second or longer")

	errTemplateInvalidPathURIs = "some path is invalid --> %+v"

	errInvalidConfirmInput = "abort because max retry count of input reached"
)

func (o *Option) ParseOptionsFromJson() error {
	if _, err := os.Stat(o.SettingJSONFilePath); err != nil {
		return fmt.Errorf("not found file: %s", o.SettingJSONFilePath)
	}

	bt, err := os.ReadFile(o.SettingJSONFilePath)
	if err != nil {
		return fmt.Errorf("cannot read file: %s", o.SettingJSONFilePath)
	}

	err = json.Unmarshal(bt, &o)
	if err != nil {
		return fmt.Errorf("bad format json: %s", o.SettingJSONFilePath)
	}

	return nil
}

func (o *Option) Valid() error {
	switch {
	case len(option.DistributionId) == 0:
		return errInvalidDistributionId
	case len(option.Paths) == 0:
		return errInvalidPathLength
	case option.RetryInterval < 1:
		return errInvalidRetryInterval
	}

	if err := option.Paths.ValidateURIs(); err != nil {
		return err
	}

	for _, paths := range o.PathGroups {
		if err := paths.ValidateURIs(); err != nil {
			return err
		}
	}

	return nil
}

func (o *Option) GetPaths() []Paths {
	paths := []Paths{}
	if len(o.Paths) > 0 {
		paths = append(paths, o.Paths)
	}

	if len(o.PathGroups) > 0 {
		paths = append(paths, o.PathGroups...)
	}

	return paths
}

func (o *Option) printExecuteContents() {
	printTemplate := `---
[ Profile        ] %s
[ DistributionId ] %s
[ Paths          ] %s
[ RetryInterval  ] %s
[ MaxRetryCount  ] %s
---
`
	var (
		profile       string
		paths         []string
		retryInterval = fmt.Sprintf("%d sec", o.RetryInterval)
		maxRetryCount string
	)
	if len(o.Profile) == 0 {
		profile = "<default profile>"
	} else {
		profile = o.Profile
	}

	for i := range o.GetPaths() {
		paths = append(paths, fmt.Sprintf("\n- %+v", o.PathGroups[i]))
	}

	if o.MaxRetryCount == 0 {
		maxRetryCount = "<infinite>"
	} else {
		maxRetryCount = strconv.Itoa(o.MaxRetryCount)
	}

	joinedPaths := strings.Join(paths, "")

	fmt.Printf(printTemplate, profile, o.DistributionId, joinedPaths, retryInterval, maxRetryCount)
}

func (o *Option) Confirm() (bool, error) {
	o.printExecuteContents()

	if o.Force {
		fmt.Println("force execute")
		return true, nil
	}

	retryCount := 0

	for {
		fmt.Print("execute? (y/n) ->")
		sc := bufio.NewScanner(os.Stdin)
		sc.Scan()
		input := sc.Text()

		switch input {
		case "y":
			return true, nil
		case "n":
			return false, nil
		}

		retryCount++
		if retryCount > 2 {
			return false, errors.New(errInvalidConfirmInput)
		}
	}
}
