package ak8sv

import (
	"context"
	ctx "context"
	"log"
	"net/url"
	"os"
	"path"
	"strings"

	"github.com/Azure/azure-sdk-for-go/profiles/2019-03-01/keyvault/keyvault"
	kvauth "github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
)

func filterSecret(s keyvault.SecretItem, fi []string, fe []string) bool {
	// This requires ALL tags be present on a secret, may need to refactor
	for _, t := range fi {
		if _, hit := s.Tags[t]; !hit {
			return false
		}
	}
	for _, t := range fe {
		if _, hit := s.Tags[t]; hit {
			return false
		}
	}
	log.Printf("%s added to secret\n", path.Base(*s.ID))
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
	sIterator, err := kv.GetSecretsComplete(ctx.Background(), GetKvURL(kvName), nil)
	if err != nil {
		log.Printf("Unable to retrieve secrets: %v", err.Error())
		os.Exit(1)
	}

	for sIterator.NotDone() {
		if filterSecret(sIterator.Value(), kvTagsInc, kvTagsEx) {
			l = append(l, path.Base(*sIterator.Value().ID))
			fCount++
		}
		err := sIterator.Next()
		if err != nil {
			log.Printf("Failed to iterator keyvault secrets: %v", err.Error())
			os.Exit(1)
		}
	}
	log.Printf("%v filtered results will be added to the secret\n", fCount)
	return l
}

// GetSecret - Retrieve the value of a KV secret
func GetSecret(s string) string {
	log.Printf("Getting secret %v....\n", s)
	v, err := kv.GetSecret(context.Background(), GetKvURL(kvName), s, "")
	if err != nil {
		log.Printf("Failed to get value for %v.\n", s)
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
