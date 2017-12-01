package action

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/rthallisey/service-broker-ci/pkg/runtime"
)

func Bind(repo string, cmd string, target string) error {
	template, err := getTemplate(repo, "template")
	if err != nil {
		return err
	}

	fmt.Printf("Running: %s create -f %s\n", cmd, template)
	args := fmt.Sprintf("create -f %s", template)
	output, err := runtime.RunCommand(cmd, args)

	fmt.Println(string(output))
	if err != nil {
		return err
	}

	// Wait for binding resource creation
	err = waitUntilReady(repo)
	if err != nil {
		return err
	}

	targetName := resourceName(target)
	instanceName, err := runtime.RunCommand("kubectl", fmt.Sprintf("get -f /tmp/%s.yaml -o jsonpath='{ .metadata.name }'", targetName))
	if err != nil {
		fmt.Println(string(instanceName))
		return err
	}
	fmt.Printf("Using Instance Name: %s\n", instanceName)

	// Get the name of the secret
	repoName := resourceName(repo)
	secretName, err := runtime.RunCommand("kubectl", fmt.Sprintf("get -f /tmp/%s -o jsonpath='{ .spec.secretName }'", repoName))
	if err != nil {
		fmt.Println(secretName)
		return err
	}

	err = waitUntilResourceReady(string(secretName), "secret")
	if err != nil {
		return err
	}

	// Gather bind data from secret
	bindData, err := runtime.RunCommand("kubectl", fmt.Sprintf("get secret %s -o jsonpath='{.data}'", secretName))
	if err != nil {
		fmt.Println(bindData)
		return err
	}
	data := strings.TrimPrefix(string(bindData), "map[")
	data = strings.TrimSuffix(data, "]")

	decodeMap := strings.Split(data, " ")
	var dataString string

	// decode base64 bindData
	for _, m := range decodeMap {
		keyValue := strings.Split(m, ":")
		decoded, err := base64.StdEncoding.DecodeString(keyValue[1])
		if err != nil {
			return err
		}
		dataString = fmt.Sprintf("%s %s=%s", dataString, keyValue[0], decoded)
	}

	fmt.Printf("Looking for a Deployment with the SAME name used in your ServiceInstance: %s\n", instanceName)

	// Inject bind data into the pod
	// oc env dc mediawiki123 DB_HOST=postgres DB_NAME=admin
	var attempt int
	retries := 60
	for attempt = 0; attempt < retries; attempt++ {
		out, err := runtime.Runtime.InjectBindData(instanceName, dataString)
		fmt.Println(string(out))
		if err == nil {
			break
		}
		attempt += 1
		time.Sleep(time.Duration(5) * time.Second)
	}
	if attempt == retries {
		return errors.New(fmt.Sprint("Timed out updating deployment %s", instanceName))
	}

	return nil
}
