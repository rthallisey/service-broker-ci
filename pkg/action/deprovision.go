package action

import (
	"fmt"
	"strings"
)

func Deprovision(repo string, cmd string) error {
	template := fmt.Sprintf("/tmp/%s", resourceName(repo))
	fmt.Printf("Running: %s delete -f %s\n", cmd, template)
	args := fmt.Sprintf("delete -f %s", template)
	output, err := RunCommand(cmd, strings.Fields(args))

	fmt.Println(string(output))
	if err != nil {
		return err
	}
	return nil
}
