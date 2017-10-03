package ci

import (
	"bufio"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

func newConfig() (*Config, error) {
	config := new(Config)
	config.ActionList = make([]map[string]string, 0)
	action, repo := "", ""

	if _, err := os.Stat(ConfigFile); err != nil {
		return nil, err
	}

	if file, err := os.Open(ConfigFile); err == nil {
		defer file.Close()

		// create a new scanner and read the file line by line
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			a := &YamlActions{}
			actionMap := make(map[string]string)

			if err = yaml.Unmarshal([]byte(scanner.Text()), a); err != nil {
				fmt.Println(err)
			}
			action, repo = getAction(a)
			if action != "" || repo != "" {
				actionMap[action] = repo
				config.ActionList = append(config.ActionList, actionMap)
			}

		}
		if err = scanner.Err(); err != nil {
			return nil, err
		}
	} else {
		return nil, err
	}
	return config, nil
}

func getAction(a *YamlActions) (string, string) {
	if a == nil {
		return "", ""
	} else if a.Provision != "" {
		return "provision", a.Provision
	} else if a.Bind != "" {
		return "bind", a.Bind
	} else if a.Unbind != "" {
		return "unbind", a.Unbind
	} else if a.Deprovision != "" {
		return "deprovision", a.Deprovision
	} else if a.Verify != "" {
		return "verify", a.Verify
	} else {
		return "", ""
	}
}
