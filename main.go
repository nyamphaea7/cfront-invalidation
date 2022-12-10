package main

import (
	"fmt"
	"os"

	"github.com/nyamphaea7/cfront-invalidation/cfront"
)

func main() {
	args := cfront.GetArgs()

	err := cfront.CreateInvalidations(args.InvalidationGroupIds)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
