package ak8sv

import (
	"context"
	ctx "context"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
)

// GetKvURL - Turn standard KV name into full URL
func GetKvURL(kvName string) string {
	// vault.azure.net
	_, err := url.Parse(kvName)
	if err == nil && strings.Contains(kvName, "vault.azure.net") {
		// Seemingly valid KV url
		return kvName
	}
	return "https://" + kvName + ".vault.azure.net"
}

// GetSecretList - Get the names of all sevret names in keyvault
func GetSecretList() []string {
	var l []string
	lResp, err := kv.GetSecrets(ctx.Background(), GetKvURL(kvName), nil)
	if err != nil {
		fmt.Println("Unable to retrieve secrets:")
		panic(err.Error())
	}
	for _, i := range lResp.Values() {
		l = append(l, path.Base(*i.ID))
	}
	return l
}

// GetSecret - Retrieve the value of a KV secret
func GetSecret(s string) string {
	v, err := kv.GetSecret(context.Background(), GetKvURL(kvName), s, "")
	if err != nil {
		fmt.Printf("Failed to get value for %v.\n", s)
		panic(err.Error())
	}
	return *v.Value
}

func newKvClient() keyvault.BaseClient {
	kvClient := keyvault.New()
	authorizer, err := kvauth.NewAuthorizerFromEnvironment()
	if err != nil {
		panic(err.Error())
	}
	kvClient.Authorizer = authorizer
	return kvClient
}
