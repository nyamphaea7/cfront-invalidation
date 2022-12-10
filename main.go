package main

import (
	"fmt"

	"github.com/nyamphaea7/cfront-invalidation/cfront"
)

func main() {
	args := cfront.GetArgs()

	fmt.Printf("ids: %v(%d)\n", args.InvalidationGroupIds, len(args.InvalidationGroupIds))
}
