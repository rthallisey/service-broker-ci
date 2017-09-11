package action

import (
	"fmt"
)

func Bind(repo string, cmd string, target string) error {
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
	// waitForResource()
	// errorCheck()
	return nil
}
