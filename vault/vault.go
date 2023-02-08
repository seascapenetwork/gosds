// Keep the credentials in a vault
package vault

import (
	"context"
	"errors"
	"fmt"

	"github.com/blocklords/gosds/env"
	hashicorp "github.com/hashicorp/vault/api"
)

type Vault struct {
	client  *hashicorp.Client
	context context.Context
	path    string
}

// Sets up the connection to the Hashicorp Vault
// If you run the Vault in the dev mode, then path should be "secret/"
func New() (*Vault, error) {
	if !env.Exists("SDS_VAULT_HOST") {
		return nil, errors.New("missing 'SDS_VAULT_HOST' environment variable")
	}
	if !env.Exists("SDS_VAULT_PORT") {
		return nil, errors.New("missing 'SDS_VAULT_PORT' environment variable")
	}
	if !env.Exists("SDS_VAULT_SECURE") {
		return nil, errors.New("missing 'SDS_VAULT_SECURE' environment variable")
	}

	if !env.Exists("SDS_VAULT_TOKEN") {
		return nil, errors.New("missing 'SDS_VAULT_TOKEN' environment variable")
	}

	secure := env.GetString("SDS_VAULT_SECURE")
	if secure != "false" && secure != "true" {
		return nil, errors.New("the value of 'SDS_VAULT_SECURE' could be 'false' or 'true'")
	}

	host := env.GetString("SDS_VAULT_HOST")
	port := env.GetString("SDS_VAULT_PORT")
	path := ""

	config := hashicorp.DefaultConfig()
	if secure == "true" {
		path = "sds"
		config.Address = fmt.Sprintf("https://%s:%s", host, port)
	} else {
		path = "secret"
		config.Address = fmt.Sprintf("http://%s:%s", host, port)
	}

	client, err := hashicorp.NewClient(config)
	if err != nil {
		return nil, err
	}

	client.SetToken(env.GetString("SDS_VAULT_TOKEN"))

	ctx := context.TODO()

	return &Vault{client: client, context: ctx, path: path}, nil
}

// Returns the String in the secret, by key
func (v *Vault) GetString(secret_name string, key string) (string, error) {
	secret, err := v.client.KVv2(v.path).Get(v.context, secret_name)
	if err != nil {
		return "", err
	}

	value, ok := secret.Data[key].(string)
	if !ok {
		fmt.Println(secret)
		return "", fmt.Errorf("vault error. failed to get the key %T %#v", secret.Data[key], secret.Data[key])
	}

	return value, nil
}
