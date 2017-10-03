package action

import (
	"fmt"
	"strings"
)

func Verify(repo string, args string) error {
	script, err := getTemplate(repo)
	if err != nil {
		return err
	}

	fmt.Printf("Running: %s %s\n", script, args)
	combinedArgs := fmt.Sprintf("%s %s", script, args)
	output, err := RunCommand("bash", strings.Fields(combinedArgs))
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found in $PATH") {
			output, err = RunCommand(script, strings.Fields(args))
			fmt.Println(string(output))
			if err != nil {
				return err
			}
		} else {
			fmt.Println(string(output))
			return err
		}
	} else {
		fmt.Println(string(output))
	}
	return nil
}
