package ci

import (
	"fmt"

	flags "github.com/jessevdk/go-flags"
)

type Args struct {
	Cluster string `short:"c" long:"cluster" description:"Cluster (openshift/kubernetes)" default:"openshift"`
}

func GetArgs() (Args, error) {
	args := Args{}

	_, err := flags.Parse(&args)
	if err != nil {
		fmt.Printf("err - %v", err)
		return args, err
	}
	return args, nil
}
