package ci

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func (c *Config) errorCheck() {

}

func (c *Config) Provision(repo string) error {
	r := strings.Split(repo, "/")

	// [ ansibleplaybookbundle, mediawiki123 ]
	gitOrg, apb := r[0], r[1]

	// For testing, set gitOrg to rthallisey
	gitOrg = "rthallisey"

	// APB template will be in the template directory
	templateAddress := fmt.Sprintf("%s/%s/%s/template/%s.yaml", BaseURL, gitOrg, Branch, apb)

	template, err := downloadTemplate(templateAddress)
	if err != nil {
		return err
	}

	output, err := RunCommand(c.Cluster, "create", "-f", template)
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	// c.waitForResource()
	// c.errorCheck()
	return nil
}

func (c *Config) Bind(repo string) error {
	return nil
}

func (c *Config) Deprovision(repo string) error {
	return nil
}

func (c *Config) Unbind(repo string) error {
	return nil
}

func (c *Config) Verify(repo string) error {
	return nil
}

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
