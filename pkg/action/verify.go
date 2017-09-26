package action

import (
	"fmt"
	"strings"
)

func Verify(repo string, args string) error {
	var script string
	var err error

	if strings.Contains(repo, "https://raw.githubusercontent.com") {
		script, err = downloadTemplate(repo)
		if err != nil {
			return err
		}
	} else {
		script = repo
	}

	fmt.Printf("Running: %s %s\n", script, args)
	output, err := RunCommand(script, strings.Fields(args))
	if err != nil {
		if strings.Contains(err.Error(), "executable file not found in $PATH") {
			combinedArgs := fmt.Sprintf("%s %s", script, args)
			output, err = RunCommand("bash", strings.Fields(combinedArgs))
		} else {
			return err
		}
	}

	fmt.Println(string(output))

	return nil
}
