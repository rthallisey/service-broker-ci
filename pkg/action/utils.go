package action

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func getTemplate(repo string, dir string) (string, error) {
	var template string
	var err error
	var args string

	if strings.Contains(repo, "https://raw.githubusercontent.com") {
		template, err = downloadTemplate(repo)
		if err != nil {
			return "", err
		}
	} else {

		if dir == "template" {
			r := strings.Split(repo, "templates/")[1]
			args = fmt.Sprintf("%s /tmp/%s", repo, r)
			template = fmt.Sprintf("/tmp/%s", r)
		} else {
			args = fmt.Sprintf("%s /tmp/%s", repo, repo)
			template = fmt.Sprintf("/tmp/%s", repo)
		}

		output, err := RunCommand("cp", args)
		if err != nil {
			// If the file doesn't exist, error later
			if !strings.Contains(string(output), "No such file or directory") {
				fmt.Println(string(output))
				return "", err
			}
		}

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
	fullCMD := append([]string{"-c"}, []string{combinedCMD}...)
	output, err := exec.Command("bash", fullCMD...).CombinedOutput()
	return output, err
}

func getObjectStatus(appName string) (string, string) {
	statusReason := fmt.Sprintf("get -f %s -o jsonpath='{ .status.conditions[0].reason }'", appName)
	reason, err := RunCommand("oc", statusReason)
	if err != nil {
		return "", ""
	}

	statusMessage := fmt.Sprintf("get -f %s -o jsonpath='{ .status.conditions[0].message }'", appName)
	message, err := RunCommand("oc", statusMessage)
	if err != nil {
		return "", ""
	}

	return string(reason), string(message)
}

const (
	Retries = 60

	Provisioned = "ProvisionedSuccessfully"
	Binded      = "InjectedBindResult"
)

func waitUntilReady(resourceName string) error {
	fmt.Printf("Waiting for %s to be ready\n", resourceName)
	var attempt int
	for attempt = 0; attempt < Retries; attempt++ {
		reason, status := getObjectStatus(resourceName)
		if reason == Provisioned || reason == Binded {
			fmt.Println("Bind or provison completed")
			fmt.Println(status)
			break
		}
		fmt.Println(reason)

		attempt += 1
		time.Sleep(time.Duration(5) * time.Second)
	}
	if attempt == Retries {
		return errors.New("Timed out waiting for resource")
	}
	return nil
}

func waitUntilResourceReady(resourceName string, resourceType string) error {
	var attempt int
	for attempt = 0; attempt < Retries; attempt++ {
		output, err := RunCommand("oc", fmt.Sprintf("get %s %s", resourceType, resourceName))

		if err == nil {
			break
		}
		fmt.Println(string(output))

		attempt += 1
		time.Sleep(time.Duration(5) * time.Second)
	}
	if attempt == Retries {
		return errors.New("Timed out waiting for resource")
	}

	return nil
}

func waitUntilDeleted(resourceName string) error {
	var attempt int
	for attempt = 0; attempt < Retries; attempt++ {
		output, err := RunCommand("oc", fmt.Sprintf("get -f %s", resourceName))

		// RunCommand errors if the resource has been deleted
		if err != nil {
			break
		}
		fmt.Println(string(output))

		attempt += 1
		time.Sleep(time.Duration(5) * time.Second)
	}
	if attempt == Retries {
		return errors.New("Timed out waiting for resource")
	}

	return nil
}
