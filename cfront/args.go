package cfront

import "flag"

type CommandLineArguments struct {
	InvalidationGroupIds InvalidationGroupIds
}

func GetArgs() CommandLineArguments {
	var ids InvalidationGroupIds
	flag.Var(&ids, "i", ids.GetArgDescription())

	flag.Parse()
	return CommandLineArguments{
		InvalidationGroupIds: ids,
	}
}
