package ci

type Config struct {
	Cluster    string
	ActionList []map[string]string

	// The bindTarget is the first available Provisioned APB that is not
	// mentioned in Bind.
	Provisioned []string
}

type YamlActions struct {
	Provision   string `yaml:"provision"`
	Bind        string `yaml:"bind"`
	Unbind      string `yaml:"unbind"`
	Deprovision string `yaml:"deprovision"`
	Verify      string `yaml:"verify"`
}

const (
	BaseURL    = "https://raw.githubusercontent.com"
	Branch     = "master"
	ConfigFile = "config.yaml"
)
