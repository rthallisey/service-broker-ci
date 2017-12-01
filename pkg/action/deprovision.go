package action

import (
	"fmt"

	"github.com/rthallisey/service-broker-ci/pkg/runtime"
)

func Deprovision(repo string, cmd string) error {
	template := fmt.Sprintf("/tmp/%s", resourceName(repo))
	fmt.Printf("Running: %s delete -f %s\n", cmd, template)
	args := fmt.Sprintf("delete -f %s", template)
	output, err := runtime.RunCommand(cmd, args)

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
