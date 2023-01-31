package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	err := GetOptions()
	if err != nil {
		fmt.Println(err)
		flag.Usage()
		os.Exit(1)
	}

	confirmExecute, err := option.Confirm()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if !confirmExecute {
		os.Exit(0)
	}

	cfg, err := MakeAwsConfig(option)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cf := CFront{
		Opt: &option,
	}

	err = cf.CreateInvalidation(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
