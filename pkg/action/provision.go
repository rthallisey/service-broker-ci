package action

import (
	"fmt"
)

func Provision(repo string, cmd string) error {
	template, err := getTemplate(repo, "template")
	if err != nil {
		return err
	}

	fmt.Printf("Running: %s create -f %s\n", cmd, template)
	args := fmt.Sprintf("create -f %s", template)
	output, err := RunCommand(cmd, args)

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
