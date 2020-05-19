package ak8sv

import (
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

func newKvClient() keyvault.BaseClient {
	kvClient := keyvault.New()
	authorizer, err := kvauth.NewAuthorizerFromEnvironment()
	if err != nil {
		panic(err.Error())
	}
	kvClient.Authorizer = authorizer
	return kvClient
}
