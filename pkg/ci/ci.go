package ci

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

func CreateCi() (*Config, error) {
	conf, err := newConfig()
	if err != nil {
		return nil, err
	}

	args, err := GetArgs()
	if err != nil {
		return nil, err
	}
	conf.setCluster(args.Cluster)

	return conf, nil
}

func (c *Config) Run() {
	for _, v := range c.ActionList {
		for action, repo := range v {
			if action != "" || repo != "" {
				fmt.Printf("ACTION: %s\n", action)
				fmt.Printf("REPO: %s\n", repo)
				err := c.callAction(action, repo)
				if err != nil {
					fmt.Println(err.Error())
					os.Exit(1)
				}
			}
		}
	}
}

// Call the action the user wants
func (c *Config) callAction(action string, repo string) error {
	var err error
	if repo == "" {
		return errors.New(fmt.Sprintf("YAML Error: Empty string used with %s action", action))
	}

	if action == "provision" {
		err = c.Provision(repo)
	} else if action == "bind" {
		err = c.Bind(repo)
	} else if action == "deprovision" {
		err = c.Deprovision(repo)
	} else if action == "unbind" {
		err = c.Unbind(repo)
	} else if action == "verify" {
		err = c.Verify(repo)
	} else {
		return errors.New(fmt.Sprintf("Action %s not found. Valid actions: [provision, bind, unbind, deprovision, verify]", action))
	}

	if err != nil {
		return err
	}

	return nil
}

// Verify the cluster we're running on
func (c *Config) setCluster(client string) error {
	client = strings.ToLower(client)
	if client == "openshift" {
		c.Cluster = "oc"
		fmt.Println("Using OpenShift Cluster")
	} else if client == "kubernetes" {
		c.Cluster = "kubectl"
		fmt.Println("Using Kubernetes Cluster")
	} else {
		fmt.Println("Using unsupported Cluster: %v", client)
		return errors.New("Unsupported Cluster")
	}
	return nil
}
