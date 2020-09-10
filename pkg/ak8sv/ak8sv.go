package ak8sv

import (
	"log"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/keyvault/keyvault"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
)

var (
	// AKV Tags can contain whitespace we don't trim that here, be careful creating this var
	kvTagsInc  []string              = strings.Split(os.Getenv("KEYVAULT_TAGS_INCLUDE"), ",")
	kvTagsEx   []string              = strings.Split(os.Getenv("KEYVAULT_TAGS_EXCLUDE"), ",")
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
	log.Printf("Keyvault:\t%v\n", kvName)
	log.Printf("URL:\t%v\n", GetKvURL(kvName))
	log.Printf("Secret:\t%v/%v\n", sNamespace, sName)
	log.Printf("Included Tags:\t%v\n", kvTagsInc)
	log.Printf("Excluded Tags:\t%v\n\n", kvTagsEx)
	switch sType {
	case "config":
		s = NewConfigSecret()
	case "certificate":
		log.Println("Certificates not implemented yet.")
	default:
		log.Println("Unsupported secret type provided, exiting")
		os.Exit(1)
	}
	ApplySecret(s)
}

// InitEnvData - Ingest environment variables to configure app
func initEnvData(e string) string {
	v := os.Getenv(e)
	if len(v) == 0 {
		log.Printf("[AK8SV] ERROR: Unable to read %v from environment!\n", e)
		os.Exit(1)
	}
	return v
}
