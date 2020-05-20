package ak8sv

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/keyvault/keyvault"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	k8s        *kubernetes.Clientset = newK8sClient()
	kv         keyvault.BaseClient   = newKvClient()
	kvName     string                = initEnvData("KEYVAULT_NAME")
	sName      string                = initEnvData("SECRET_NAME")
	sNamespace string                = initEnvData("SECRET_NAMESPACE")
	sType      string                = initEnvData("SECRET_TYPE")
)

// Bootstrap - Entry point for the application
func Bootstrap() {
	var s apiv1.Secret
	// General config dumping for easier debugging
	fmt.Printf("Using Keyvault: %v\n", kvName)
	fmt.Printf("URL: %v\n", GetKvURL(kvName))
	fmt.Printf("Secret: %v/%v\n", sNamespace, sName)
	switch sType {
	case "config":
		s = NewConfigSecret()
	case "certificate":
		fmt.Println("Not implemented yet.")
	default:
		fmt.Println("Unsupported secret type provided, exiting.")
		os.Exit(1)
	}
	ApplySecret(s)
	fmt.Println("Secret updated successfully.")
}

// InitEnvData - Ingest environment variables to configure app
func initEnvData(e string) string {
	v := os.Getenv(e)
	if len(v) == 0 {
		fmt.Printf("ERROR: Unable to read %v from environment!\n", e)
		os.Exit(1)
	}
	return v
}
