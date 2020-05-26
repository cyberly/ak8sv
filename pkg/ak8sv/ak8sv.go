package ak8sv

import (
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/keyvault/keyvault"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	kvInclude  string
	kvExclude  string
	kv         keyvault.BaseClient   = newKvClient()
	kvName     string                = initEnvData("KEYVAULT_NAME")
	k8s        *kubernetes.Clientset = newK8sClient()
	sName      string                = initEnvData("SECRET_NAME")
	sNamespace string                = initEnvData("SECRET_NAMESPACE")
	sType      string                = initEnvData("SECRET_TYPE")
)

// Bootstrap - Entry point for the application
func Bootstrap() {
	var s apiv1.Secret
	fmt.Println("AK8sV Config:")
	fmt.Println("Keyvault:")
	fmt.Printf("Keyvault:\t%v\n", kvName)
	fmt.Printf("\t\t%v\n", GetKvURL(kvName))
	fmt.Printf("Secret:\t%v/%v\n\n", sNamespace, sName)
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
