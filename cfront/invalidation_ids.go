package cfront

import (
	"fmt"
	"strings"
)

type InvalidationGroupIds []string

func (ss *InvalidationGroupIds) String() string {
	return fmt.Sprintf("%s", *ss)
}

func (ss *InvalidationGroupIds) Set(value string) error {
	split := strings.Split(value, ",")
	*ss = append(*ss, split...)
	return nil
}

func (InvalidationGroupIds) GetArgDescription() string {
	return strings.Join([]string{
		"Execute only invalidation group ids from settings.json.",
		"e.q. go run ./main.go -i XX,YY",
		"e.q. go run ./main.go -i XX -i YY",
	}, "\n")
}
