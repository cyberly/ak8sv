package ak8sv

import (
	"fmt"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/keyvault/keyvault"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	// AKV Tags can contain whitespace we don't trim that here, be careful creating this var
	kvTagsInc  []string              = strings.Split(os.Getenv("KEYVAULT_TAGS_EXCLUDE"), ",")
	kvTagsEx   []string              = strings.Split(os.Getenv("KEYVAULT_TAGS_INCLUDE"), ",")
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
	fmt.Println("[AK8SV] Config:")
	fmt.Println("[AK8SV] Keyvault:")
	fmt.Printf("[AK8SV] Keyvault:\t%v\n", kvName)
	fmt.Printf("[AK8SV]\t\t%v\n", GetKvURL(kvName))
	fmt.Printf("[AK8SV] Secret:\t%v/%v\n\n", sNamespace, sName)
	switch sType {
	case "config":
		if len(kvTagsInc) > 0 && len(kvTagsEx) > 0 {
			fmt.Println("[AK8SV] WARNING: Excluded tags will superceded included tags.")
		}
		s = NewConfigSecret()
	case "certificate":
		fmt.Println("[AK8SV] Not implemented yet.")
	default:
		fmt.Println("[AK8SV] Unsupported secret type provided, exiting")
		os.Exit(1)
	}
	ApplySecret(s)
}

// InitEnvData - Ingest environment variables to configure app
func initEnvData(e string) string {
	v := os.Getenv(e)
	if len(v) == 0 {
		fmt.Printf("[AK8SV] ERROR: Unable to read %v from environment!\n", e)
		os.Exit(1)
	}
	return v
}
