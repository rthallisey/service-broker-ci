package ci

import (
	"errors"
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
	target := strings.Split(bindTarget, "/")[1]

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
	apb := r[len(r)-1]
	gitOrg := strings.Join(r[0:len(r)-1], "/")

	// APB template will be in the template directory
	return fmt.Sprintf("%s/%s/%s/template/%s.yaml", BaseURL, gitOrg, Branch, apb)
}

func findBindTarget(repo string, provisioned []string) (string, []string, error) {
	var usedTargets []int
	foundTarget := false
	foundBind := false
	var bindTarget string

	// The config in imperative so order matters
	for count, r := range provisioned {
		// Remove the first Provisioned app that matches the Bind repo
		// and the first Provisioned app that doesn't.

		// The first Provisioned app that doesn't match Bind is the
		// bindTarget.
		if r != repo && !foundTarget {
			bindTarget = r
			foundTarget = true
			usedTargets = append(usedTargets, count)
		}

		// The first Provisioned app that matches the Bind repo is the
		// bind app.
		if r == repo && !foundBind {
			foundBind = true
			usedTargets = append(usedTargets, count)
		}

		if foundBind && foundTarget {
			cleanupUsedTargets(usedTargets, provisioned)
			return bindTarget, provisioned, nil
		}
	}

	return "", provisioned, errors.New("Failed to find a provisioned bind target and bind app")
}

func cleanupUsedTargets(usedTargets []int, provisioned []string) []string {
	for _, d := range usedTargets {
		// Cleanup the bindTarget and the bind app
		if len(provisioned) == 1 {
			provisioned = provisioned[:0]
		} else {
			provisioned = append(provisioned[:d], provisioned[d+1:]...)
		}
	}
	return provisioned
}
