package action

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func getTemplate(repo string) (string, error) {
	var template string
	var err error

	if strings.Contains(repo, "https://raw.githubusercontent.com") {
		template, err = downloadTemplate(repo)
		if err != nil {
			return "", err
		}
	} else {
		args := fmt.Sprintf("%s /tmp/%s", repo, repo)
		output, err := RunCommand("cp", args)
		if err != nil {
			// If the file doesn't exist, error later
			if !strings.Contains(string(output), "No such file or directory") {
				fmt.Println(string(output))
				return "", err
			}
		}

		template = fmt.Sprintf("/tmp/%s", repo)
	}

	return template, nil
}

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

func RunCommand(cmd string, args string) ([]byte, error) {
	combinedCMD := fmt.Sprintf("%s %s", cmd, args)
	// a := append([]string{cmd}, args...)
	// t := strings.Join(a, " ")
	fullCMD := append([]string{"-c"}, []string{combinedCMD}...)
	output, err := exec.Command("bash", fullCMD...).CombinedOutput()
	return output, err
}
