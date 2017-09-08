package action

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func downloadTemplate(url string) (string, error) {
	p := strings.Split(url, "/")
	path := fmt.Sprintf("/tmp/%s", p[len(p)-1])
	out, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer out.Close()

	req, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer req.Body.Close()

	_, err = io.Copy(out, req.Body)
	if err != nil {
		return "", err
	}

	return path, nil
}

func RunCommand(cmd string, args ...string) ([]byte, error) {
	output, err := exec.Command(cmd, args...).CombinedOutput()
	return output, err
}
