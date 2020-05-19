package ak8sv

import (
	ctx "context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	sData                            = make(map[string][]byte)
	sName      string                = InitEnvData("SECRET_NAME")
	sNamespace string                = InitEnvData("SECRET_NAMESPACE")
	sType      string                = InitEnvData("SECRET_TYPE")
	k8s        *kubernetes.Clientset = newK8sClient()
)

// ApplySecret - Apply the secret configured by ingested environment variables
func ApplySecret(data apiv1.Secret) *apiv1.Secret {
	var secretResp *apiv1.Secret
	var err error
	if checkSecret() {
		fmt.Println("Updating secret...")
		secretResp, err = k8s.CoreV1().Secrets(sNamespace).Update(ctx.TODO(), &data, metav1.UpdateOptions{})
	} else {
		fmt.Println("Updating secret...")
		secretResp, err = k8s.CoreV1().Secrets(sNamespace).Create(ctx.TODO(), &data, metav1.CreateOptions{})
	}
	if err != nil {
		fmt.Printf("Error: %v\n", err.Error())
		fmt.Println("Applying secret failed, exiting.")
		os.Exit(1)
	}
	fmt.Println("Secret applied successfully.")
	return secretResp
}

func checkSecret() bool {
	_, err := k8s.CoreV1().Secrets(sNamespace).Get(ctx.TODO(), sName, metav1.GetOptions{})
	if err != nil {
		return false
	}
	return true
}

func newK8sClient() *kubernetes.Clientset {
	config, err := rest.InClusterConfig()
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func newK8sClientLocal() kubernetes.Clientset {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String(
			"kubeconfig", filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return *clientset
}

func newConfigSecret(sData map[string][]byte) apiv1.Secret {
	s := apiv1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      sName,
			Namespace: sNamespace,
		},
		Data: sData,
		Type: "Opaque",
	}
	return s
}
