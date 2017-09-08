package action

import (
	"fmt"
)

func Unbind(repo string, cmd string) error {
	fmt.Printf("Running: %s delete ServiceInstanceCredential binding\n", cmd)
	output, err := RunCommand(cmd, "delete", "ServiceInstanceCredential", "binding")
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
