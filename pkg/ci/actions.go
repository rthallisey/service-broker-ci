package ci

import (
	"fmt"
	"strings"

	"github.com/rthallisey/service-broker-ci/pkg/action"
)

type Broker interface {
	Provision(string) error
	Deprovision(string) error
	Bind(string) error
	Unbind(string) error
	Verify(string) error
}

func (c *Config) Provision(repo string) error {
	err := action.Provision(getTemplateAddr(repo), c.Cluster)
	if err != nil {
		return err
	}

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

func getTemplateAddr(repo string) string {
	r := strings.Split(repo, "/")

	// [ ansibleplaybookbundle, mediawiki123 ]
	gitOrg, apb := r[0], r[1]

	// For testing, set gitOrg to rthallisey/service-broker-ci
	gitOrg = "rthallisey/service-broker-ci"

	// APB template will be in the template directory
	return fmt.Sprintf("%s/%s/%s/template/%s.yaml", BaseURL, gitOrg, Branch, apb)
}
