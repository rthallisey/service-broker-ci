package runtime

import (
	"fmt"
	"os/exec"
)

func RunCommand(cmd string, args string) ([]byte, error) {
	combinedCMD := fmt.Sprintf("%s %s", cmd, args)
	fullCMD := append([]string{"-c"}, []string{combinedCMD}...)
	output, err := exec.Command("bash", fullCMD...).CombinedOutput()
	return output, err
}
