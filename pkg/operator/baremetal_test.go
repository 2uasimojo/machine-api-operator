package operator

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	fakedynamic "k8s.io/client-go/dynamic/fake"
	fakekube "k8s.io/client-go/kubernetes/fake"
	"sigs.k8s.io/yaml"
)

var yamlContent = `
apiVersion: metal3.io/v1alpha1
kind: Provisioning
metadata:
  name: test
spec:
  provisioningInterface: "ensp0"
  provisioningIP: "172.30.20.3"
  provisioningNetworkCIDR: "172.30.20.0/24"
  provisioningDHCPExternal: false
  provisioningDHCPRange: "172.30.20.10, 72.30.20.100"
`
var (
	expectedProvisioningInterface    = "ensp0"
	expectedProvisioningIP           = "172.30.20.3"
	expectedProvisioningNetworkCIDR  = "172.30.20.0/24"
	expectedProvisioningDHCPExternal = false
	expectedProvisioningDHCPRange    = "172.30.20.10, 72.30.20.100"
)

func TestGenerateRandomPassword(t *testing.T) {
	pwd := generateRandomPassword()
	if pwd == "" {
		t.Errorf("Expected a valid string but got null")
	}
}

func newOperatorWithBaremetalConfig() *OperatorConfig {
	return &OperatorConfig{
		targetNamespace,
		Controllers{
			"docker.io/openshift/origin-aws-machine-controllers:v4.0.0",
			"docker.io/openshift/origin-machine-api-operator:v4.0.0",
			"docker.io/openshift/origin-machine-api-operator:v4.0.0",
		},
		BaremetalControllers{
			"quay.io/openshift/origin-baremetal-operator:v4.2.0",
			"quay.io/openshift/origin-ironic:v4.2.0",
			"quay.io/openshift/origin-ironic-inspector:v4.2.0",
			"quay.io/openshift/origin-ironic-ipa-downloader:v4.2.0",
			"quay.io/openshift/origin-ironic-machine-os-downloader:v4.2.0",
			"quay.io/openshift/origin-ironic-static-ip-manager:v4.2.0",
		},
	}
}

//Testing the case where the password does already exist
func TestCreateMariadbPasswordSecret(t *testing.T) {
	kubeClient := fakekube.NewSimpleClientset(nil...)
	operatorConfig := newOperatorWithBaremetalConfig()
	client := kubeClient.CoreV1()

	// First create a mariadb password secret
	if err := createMariadbPasswordSecret(kubeClient.CoreV1(), operatorConfig); err != nil {
		t.Fatalf("Failed to create first Mariadb password. %s ", err)
	}
	// Read and get Mariadb password from Secret just created.
	oldMaridbPassword, err := client.Secrets(operatorConfig.TargetNamespace).Get(baremetalSecretName, metav1.GetOptions{})
	if err != nil {
		t.Fatal("Failure getting the first Mariadb password that just got created.")
	}
	oldPassword, ok := oldMaridbPassword.StringData[baremetalSecretKey]
	if !ok || oldPassword == "" {
		t.Fatal("Failure reading first Mariadb password from Secret.")
	}

	// The pasword definitely exists. Try creating again.
	if err := createMariadbPasswordSecret(kubeClient.CoreV1(), operatorConfig); err != nil {
		t.Fatal("Failure creating second Mariadb password.")
	}
	newMaridbPassword, err := client.Secrets(operatorConfig.TargetNamespace).Get(baremetalSecretName, metav1.GetOptions{})
	if err != nil {
		t.Fatal("Failure getting the second Mariadb password.")
	}
	newPassword, ok := newMaridbPassword.StringData[baremetalSecretKey]
	if !ok || newPassword == "" {
		t.Fatal("Failure reading second Mariadb password from Secret.")
	}
	if oldPassword != newPassword {
		t.Fatalf("Both passwords do not match.")
	} else {
		t.Logf("First Mariadb password is being preserved over re-creation as expected.")
	}
}

func TestGetBaremetalProvisioningConfig(t *testing.T) {
	u := &unstructured.Unstructured{Object: map[string]interface{}{}}
	if err := yaml.Unmarshal([]byte(yamlContent), &u); err != nil {
		t.Errorf("failed to unmarshall input yaml content:%v", err)
	}
	dynamicClient := fakedynamic.NewSimpleDynamicClient(runtime.NewScheme(), u)
	baremetalConfig, err := getBaremetalProvisioningConfig(dynamicClient, "test")
	if err != nil {
		t.Logf("Unstructed Config:  %+v", u)
		t.Fatalf("Failed to get Baremetal Provisioning Interface from CR %s", "test")
	}
	if baremetalConfig.ProvisioningInterface != expectedProvisioningInterface ||
		baremetalConfig.ProvisioningIp != expectedProvisioningIP ||
		baremetalConfig.ProvisioningNetworkCIDR != expectedProvisioningNetworkCIDR ||
		baremetalConfig.ProvisioningDHCPExternal != expectedProvisioningDHCPExternal ||
		baremetalConfig.ProvisioningDHCPRange != expectedProvisioningDHCPRange {
		t.Logf("Actual BaremetalProvisioningConfig: %+v", baremetalConfig)
		t.Logf("Expected : ProvisioningInterface: %s, ProvisioningIP: %s, ProvisioningNetworkCIDR: %s, ProvisioningDHCPExternal: %t, expectedProvisioningDHCPRange: %s", expectedProvisioningInterface, expectedProvisioningIP, expectedProvisioningNetworkCIDR, expectedProvisioningDHCPExternal, expectedProvisioningDHCPRange)
		t.Fatalf("failed getBaremetalProvisioningConfig. One or more BaremetalProvisioningConfig items do not match the expected config.")
	}
}

func TestGetIncorrectBaremetalProvisioningCR(t *testing.T) {
	u := &unstructured.Unstructured{Object: map[string]interface{}{}}
	if err := yaml.Unmarshal([]byte(yamlContent), &u); err != nil {
		t.Errorf("failed to unmarshall input yaml content:%v", err)
	}
	dynamicClient := fakedynamic.NewSimpleDynamicClient(runtime.NewScheme(), u)
	baremetalConfig, err := getBaremetalProvisioningConfig(dynamicClient, "test1")
	if err != nil {
		t.Logf("Unable to get Baremetal Provisioning Config from CR %s as expected", "test1")
	}
	if baremetalConfig.ProvisioningInterface != "" {
		t.Errorf("BaremetalProvisioningConfig is not expected to be set.")
	}
}
