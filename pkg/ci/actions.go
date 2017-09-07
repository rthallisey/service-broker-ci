package ci

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func (c *Config) errorCheck() {

}

func (c *Config) Provision(repo string) error {
	// req, err := http.Get(fmt.Sprintf("http://%s/apb.yaml", repo))
	req, err := http.Get("https://raw.githubusercontent.com/rthallisey/service-broker-ci/master/templates/mediawiki123.yaml")
	if err != nil {
		fmt.Println(err.Error())
	}
	defer req.Body.Close()

	data, err := ioutil.ReadAll(req.Body)

	output, err := RunCommand(c.Cluster, "create", "-f", string(data))
	if err != nil {
		return err
	}

	fmt.Println(string(output))
	return nil
	// c.waitForResource()
	// c.errorCheck()
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
