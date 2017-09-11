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

	c.Provisioned = append(c.Provisioned, repo)
	return nil
}

func (c *Config) Bind(repo string) error {
	bindTarget, p, err := findBindTarget(repo, c.Provisioned)
	if err != nil {
		return err
	}

	// Save the updated list of Provisioned apps
	c.Provisioned = p

	// Split the app name from the gitOrg
	t := strings.Split(bindTarget, "/")
	target := t[len(t)-1]

	// ansibleplaybookbundle/postgresql -> ansibleplaybookbundle/postgresql-mediawiki-bind
	//                                    <gitOrg>/<bindApp>-<bindTarget>-bind
	repo = fmt.Sprintf("%s-%s-bind", repo, target)
	err = action.Bind(getTemplateAddr(repo), c.Cluster, bindTarget)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Deprovision(repo string) error {
	err := action.Deprovision(repo, c.Cluster)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Unbind(repo string) error {
	err := action.Unbind(repo, c.Cluster)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Verify(repoAndArgs string) error {
	repo, args := getScriptAddr(repoAndArgs)
	err := action.Verify(repo, args)
	if err != nil {
		return err
	}
	return nil
}
