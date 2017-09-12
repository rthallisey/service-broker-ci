package action

import (
	"fmt"
	"strings"
)

func Unbind(repo string, cmd string) error {
	fmt.Printf("Running: %s delete ServiceInstanceCredential binding\n", cmd)
	args := fmt.Sprintf("delete ServiceInstanceCredential binding")
	output, err := RunCommand(cmd, strings.Fields(args))
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
}
