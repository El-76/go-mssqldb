//go:build go1.18
// +build go1.18

package akvkeys

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"net/url"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
)

// GetProviderCredential retrieves the Azure credential for accessing Azure Key Vault
func GetProviderCredential() (cred azcore.TokenCredential, err error) {
	sc := os.Getenv("AZURESUBSCRIPTION_SERVICE_CONNECTION_ID")
	if len(sc) > 0 {
		tenant := os.Getenv("AZURESUBSCRIPTION_TENANT_ID")
		clientID := os.Getenv("AZURESUBSCRIPTION_CLIENT_ID")
		token := os.Getenv("SYSTEM_ACCESSTOKEN")
		cred, err = azidentity.NewAzurePipelinesCredential(tenant, clientID, sc, token, nil)
		if err != nil {
			return
		}
	} else {
		cred, err = azidentity.NewDefaultAzureCredential(nil)
		if err != nil {
			return
		}
	}
	return
}

// GetTestAKV retrieves an Azure Key Vault client for testing purposes.
func GetTestAKV() (client *azkeys.Client, u string, err error) {
	vaultName := os.Getenv("KEY_VAULT_NAME")
	if len(vaultName) == 0 {
		err = fmt.Errorf("KEY_VAULT_NAME is not set in the environment")
		return
	}
	vaultURL := fmt.Sprintf("https://%s.vault.azure.net/", url.PathEscape(vaultName))
	cred, err := GetProviderCredential()
	if err != nil {
		return
	}
	client, err = azkeys.NewClient(vaultURL, cred, nil)
	if err != nil {
		return
	}
	u = vaultURL + "keys"
	return
}

func CreateRSAKey(client *azkeys.Client) (name string, err error) {
	kt := azkeys.KeyTypeRSA
	ks := int32(2048)
	rsaKeyParams := azkeys.CreateKeyParameters{
		Kty:     &kt,
		KeySize: &ks,
	}

	i, _ := rand.Int(rand.Reader, big.NewInt(1000000))
	name = fmt.Sprintf("go-mssqlkey%d", i)
	_, err = client.CreateKey(context.TODO(), name, rsaKeyParams, nil)
	if err != nil {
		_, err = client.RecoverDeletedKey(context.TODO(), name, &azkeys.RecoverDeletedKeyOptions{})
	}
	return
}

func DeleteRSAKey(client *azkeys.Client, name string) bool {
	_, err := client.DeleteKey(context.TODO(), name, nil)
	return err == nil
}
