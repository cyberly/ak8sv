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

func filterSecret(s keyvault.SecretItem, fi []string, fe []string) bool {
	for _, t := range fi {
		if _, hit := s.Tags[t]; !hit {
			return false
		}
	}
	for _, t := range fe {
		if _, hit := s.Tags[t]; !hit {
			return false
		}
	}
	fmt.Printf("[KEYVAULT] Secret %s (Inc: %v Ex: %v) included", path.Base(*s.ID), fi, fe)
	return true
}

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
	var fCount int = 0
	fmt.Println("[KEYVAULT] Retrieving secret list...")
	lResp, err := kv.GetSecrets(ctx.Background(), GetKvURL(kvName), nil)
	fmt.Printf("[KEYVAULT] Got %v secrets from key vault\n", len(lResp.Values()))
	if err != nil {
		fmt.Println("[KEYVAULT] Unable to retrieve secrets:")
		panic(err.Error())
	}
	for _, i := range lResp.Values() {
		if filterSecret(i, kvTagsInc, kvTagsEx) {
			l = append(l, path.Base(*i.ID))
			fCount++

		}
	}
	fmt.Printf("[KEYVAULT] Got %v filtered secrets will be added to the secret\n", fCount)
	return l
}

// GetSecret - Retrieve the value of a KV secret
func GetSecret(s string) string {
	fmt.Printf("[KEYVAULT] Getting secret %v....\n", s)
	v, err := kv.GetSecret(context.Background(), GetKvURL(kvName), s, "")
	if err != nil {
		fmt.Printf("[KEYVAULT] Failed to get value for %v.\n", s)
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
