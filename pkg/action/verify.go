package action

import (
	"fmt"
)

func Verify(repo string, args string) error {
	script, err := downloadTemplate(repo)
	if err != nil {
		return err
	}

	output, err := RunCommand("bash", script, args)
	if err != nil {
		return err
	}

	fmt.Println(string(output))

	return nil
}
