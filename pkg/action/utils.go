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
	fmt.Printf("URL: %s\n", url)
	path := resourceName(url)
	path = fmt.Sprintf("/tmp/%s", path)

	// Delete if it exists, so don't catch error
	os.Remove(path)

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

func resourceName(repo string) string {
	path := strings.Split(repo, "/")
	return path[len(path)-1]
}

func RunCommand(cmd string, args []string) ([]byte, error) {
	output, err := exec.Command(cmd, args...).CombinedOutput()
	return output, err
}
