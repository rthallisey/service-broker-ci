package ci

import (
	"errors"
	"fmt"
	"strings"
)

type Config struct {
	Cluster    string
	ActionList []map[string]string
}

type Actions struct {
	Provision   string `yaml:"provision"`
	Bind        string `yaml:"bind"`
	Unbind      string `yaml:"unbind"`
	Deprovision string `yaml:"deprovision"`
	Verify      string `yaml:"verify"`
}

const (
	ConfigFile = "config.yaml"
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
			}
		}
	}
}

// Verify the cluster we're running on
func (c *Config) setCluster(client string) error {
	client = strings.ToLower(client)
	if client == "openshift" {
		fmt.Println("Using OpenShift Cluster")
	} else if client == "kubernetes" {
		fmt.Println("Using Kubernetes Cluster")
	} else {
		fmt.Println("Using unsupported Cluster: %v", client)
		return errors.New("Unsupported Cluster")
	}
	c.Cluster = client
	return nil
}
