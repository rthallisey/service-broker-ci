package action

import (
	"fmt"

	"github.com/rthallisey/service-broker-ci/pkg/runtime"
)

func Provision(repo string, cmd string) error {
	template, err := getTemplate(repo, "template")
	if err != nil {
		return err
	}

	fmt.Printf("Running: %s create -f %s\n", cmd, template)
	args := fmt.Sprintf("create -f %s", template)
	output, err := runtime.RunCommand(cmd, args)

	fmt.Println(string(output))
	if err != nil {
		return err
	}

	err = waitUntilReady(template)
	if err != nil {
		return err
	}
	return nil
}
