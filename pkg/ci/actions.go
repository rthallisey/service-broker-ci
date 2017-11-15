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
	validRepo := resolveGitRepo(repo)
	r, _ := getScriptAddr(repo, validRepo, "template")
	if r == "" {
		return errors.New("Can't using an empty address for provision")
	}

	err := action.Provision(r, c.Cluster)
	if err != nil {
		return err
	}

	c.Provisioned = append(c.Provisioned, repo)
	return nil
}

func (c *Config) Bind(repo string) error {
	if repo == "" {
		return errors.New("Can't using an empty address for bind")
	}

	bindTarget, p, err := findBindTarget(repo, c.Provisioned)
	if err != nil {
		return err
	}

	// Save the updated list of Provisioned apps
	c.Provisioned = p

	// Split the app name from the gitOrg
	t := strings.Split(bindTarget, "/")
	target := t[len(t)-1]

	validRepo := resolveGitRepo(repo)

	// ansibleplaybookbundle/postgresql -> ansibleplaybookbundle/postgresql-mediawiki-bind
	//                                    <gitOrg>/<bindApp>-<bindTarget>-bind
	repo = fmt.Sprintf("%s-%s-bind", repo, target)

	r := strings.Split(repo, "/")
	resource := r[len(r)-1]
	rawGitURL := ""
	if validRepo == "" {
		rawGitURL = fmt.Sprintf("templates/%s.yaml", resource)
	} else {
		rawGitURL = fmt.Sprintf("%s/%s/%s/templates/%s.yaml", BaseURL, validRepo, Branch, resource)
	}
	fmt.Println(rawGitURL)

	err = action.Bind(rawGitURL, c.Cluster, bindTarget)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Deprovision(repo string) error {
	validRepo := resolveGitRepo(repo)
	r, _ := getScriptAddr(repo, validRepo, "template")
	if r == "" {
		return errors.New("Can't using an empty address for deprovision")
	}

	err := action.Deprovision(r, c.Cluster)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Unbind(bindInfo string) error {
	trim := strings.Replace(bindInfo, " ", "", -1)
	binding := strings.Split(trim, "|")
	err := action.Unbind(binding, c.Cluster)
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Verify(repoAndArgs string) error {
	validRepo := resolveGitRepo(repoAndArgs)
	repo, args := getScriptAddr(repoAndArgs, validRepo, "script")
	if repo == "" {
		return errors.New("Can't using an empty address for verify")
	}

	err := action.Verify(repo, args)
	if err != nil {
		return err
	}
	return nil
}
