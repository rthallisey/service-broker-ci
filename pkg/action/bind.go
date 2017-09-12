package action

import (
	"fmt"
	"strings"
)

func Bind(repo string, cmd string, target string) error {
	template, err := downloadTemplate(repo)
	if err != nil {
		return err
	}

	fmt.Printf("Running: %s create -f %s\n", cmd, template)
	args := fmt.Sprintf("create -f %s", template)
	output, err := RunCommand(cmd, strings.Fields(args))
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	// waitForResource()
	// errorCheck()
	return nil
}
