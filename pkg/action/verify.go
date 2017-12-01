package action

import (
	"fmt"
	"strings"

	"github.com/rthallisey/service-broker-ci/pkg/runtime"
)

func Verify(repo string, args string) error {
	script, err := getTemplate(repo, "script")
	if err != nil {
		return err
	}

	fmt.Printf("Running: %s %s\n", repo, args)
	combinedArgs := fmt.Sprintf("%s %s", script, args)
	output, err := runtime.RunCommand("bash", combinedArgs)
	if err != nil {
		if strings.Contains(string(output), "executable file not found in $PATH") ||
			strings.Contains(string(output), "No such file or directory") {
			output, err = runtime.RunCommand(repo, args)
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
