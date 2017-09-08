package action

import (
	"fmt"
)

func Deprovision(repo string, cmd string) error {
	resource := resourceName(repo)
	fmt.Printf("Running: %s delete ServiceInstance %s -n default\n", cmd, resource)
	output, err := RunCommand(cmd, "delete", "ServiceInstance", resource, "-n", "default")
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
