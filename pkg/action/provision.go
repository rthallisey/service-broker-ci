package action

import (
	"fmt"
	"strings"
)

func Provision(repo string, cmd string) error {
	var template string
	var err error

	if strings.Contains(repo, "https://raw.githubusercontent.com") {
		template, err = downloadTemplate(repo)
		if err != nil {
			return err
		}
	} else {
		template = repo
	}

	fmt.Printf("Running: %s create -f %s\n", cmd, template)
	args := fmt.Sprintf("create -f %s", template)
	output, err := RunCommand(cmd, strings.Fields(args))

	fmt.Println(string(output))
	if err != nil {
		return err
	}
	return nil
}
