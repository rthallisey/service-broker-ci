package action

import (
	"fmt"
	"strings"
)

func Verify(repo string, args string) error {
	script, err := downloadTemplate(repo)
	if err != nil {
		return err
	}

	fmt.Printf("Running: %s %s\n", script, args)
	output, err := RunCommand(script, strings.Fields(args))
	if err != nil {
		return err
	}

	fmt.Println(string(output))

	return nil
}
