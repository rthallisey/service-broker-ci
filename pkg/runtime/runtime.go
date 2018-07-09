package runtime

import (
	"fmt"
)

var Runtime runtime

type runtime interface {
	InjectBindData([]byte, string) ([]byte, error)
}

type openshift struct{}
type kubernetes struct{}

func GetRuntime(runtime string) {
	if runtime == "openshift" {
		Runtime = openshift{}
	} else if runtime == "kubernetes" {
		Runtime = kubernetes{}
	}
}

// InjectBindData - inject bind data using openshift
func (o openshift) InjectBindData(instanceName []byte, dataString string) ([]byte, error) {
	fmt.Printf("Looking for a Deployment Config with the SAME name used in your ServiceInstance: %s\n", instanceName)
	output, err := RunCommand("oc", fmt.Sprintf("set env dc %s %s", instanceName, dataString))
	return output, err
}

// InjectBindData - inject bind data using kubernetes
func (k kubernetes) InjectBindData(instanceName []byte, dataString string) ([]byte, error) {
	fmt.Printf("Looking for a Deployment with the SAME name used in your ServiceInstance: %s\n", instanceName)
	output, err := RunCommand("kubectl", fmt.Sprintf("set env deploy/%s %s", instanceName, dataString))
	return output, err
}
