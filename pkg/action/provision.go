package action

import (
	"fmt"
)

func Provision(repo string, cmd string) error {
	template, err := downloadTemplate(repo)
	if err != nil {
		return err
	}

	fmt.Printf("Running: %s create -f %s\n", cmd, template)
	output, err := RunCommand(cmd, "create", "-f", template)
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
