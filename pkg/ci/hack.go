package ci

import (
	"os/exec"
)

func RunCommand(cmd string, args ...string) ([]byte, error) {
	output, err := exec.Command(cmd, args...).CombinedOutput()
	return output, err
}
