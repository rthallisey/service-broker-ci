package action

import (
	"fmt"
	"strings"
)

func Deprovision(repo string, cmd string) error {
	resource := resourceName(repo)
	fmt.Printf("Running: %s delete ServiceInstance %s\n", cmd, resource)
	args := fmt.Sprintf("delete ServiceInstance %s", resource)
	output, err := RunCommand(cmd, strings.Fields(args))
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
