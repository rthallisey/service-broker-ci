package action

import (
	"fmt"
	"strings"
)

func Unbind(binding []string, cmd string) error {
	template := fmt.Sprintf("/tmp/%s-bind.yaml", strings.Join(binding, "-"))
	fmt.Printf("Running: %s delete -f %s\n", cmd, template)
	args := fmt.Sprintf("delete -f %s", template)
	output, err := RunCommand(cmd, args)

	fmt.Println(string(output))
	if err != nil {
		return err
	}

	err = waitUntilDeleted(template)
	if err != nil {
		return err
	}

	return nil
}
