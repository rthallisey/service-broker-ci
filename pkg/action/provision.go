package action

import (
	"fmt"
)

func Provision(addr string, cmd string) error {

	template, err := downloadTemplate(addr)
	if err != nil {
		return err
	}

	output, err := RunCommand(cmd, "create", "-f", template)
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	// c.waitForResource()
	// c.errorCheck()
	return nil
}
