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
	r, _ := getScriptAddr(repo, "template")
	if r == "" {
		return errors.New("Can't using an empty address for provision")
	}

	err := action.Provision(r, c.Cluster)
	if err != nil {
		return err
	}

	c.Provisioned = append(c.Provisioned, repo)

	fmt.Printf("Waiting for %s apb to be ready", action.ResourceName(repo))
	err = c.Verify(fmt.Sprintf("wait-for-resource.sh create pod %s", action.ResourceName(repo)))
	if err != nil {
		return err
	}
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

	r, _ := getScriptAddr(repo, "template")
	if r == "" {
		return errors.New("Can't using an empty address for bind")
	}

	err = action.Bind(r, c.Cluster, bindTarget)
	if err != nil {
		return err
	}

	fmt.Printf("Waiting for %s to bind to %s", repo, bindTarget)
	err = c.Verify("wait-for-resource.sh create bindings.v1alpha1.servicecatalog.k8s.io binding")
	if err != nil {
		return err
	}

	fmt.Printf("Waiting for podpreset to be created")
	err = c.Verify("wait-for-resource.sh create podpreset binding")
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Deprovision(repo string) error {
	r, _ := getScriptAddr(repo, "template")
	if r == "" {
		return errors.New("Can't using an empty address for deprovision")
	}

	err := action.Deprovision(r, c.Cluster)
	if err != nil {
		return err
	}

	// TODO: binding name needs to be a param
	fmt.Printf("Waiting for %s to be deleted", action.ResourceName(repo))
	err = c.Verify(fmt.Sprintf("wait-for-resource.sh delete pod %s", action.ResourceName(repo)))
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

	fmt.Printf("Waiting for unbind occur")
	err = c.Verify("wait-for-resource.sh delete bindings.v1alpha1.servicecatalog.k8s.io binding")
	if err != nil {
		return err
	}
	return nil
}

func (c *Config) Verify(repoAndArgs string) error {
	repo, args := getScriptAddr(repoAndArgs, "script")
	if repo == "" {
		return errors.New("Can't using an empty address for verify")
	}

	err := action.Verify(repo, args)
	if err != nil {
		return err
	}
	return nil
}
