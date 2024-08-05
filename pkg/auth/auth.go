// Package auth provides helpers for working with OCI registry authentication.
package auth

import (
	"strings"

	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	acr "github.com/chrismellard/docker-credential-acr-env/pkg/credhelper"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/v1/google"
)

// GetKeychain returns the default authn keychain derived from an OCI host.
func GetKeychain(host string) []authn.Keychain {
	keychain := []authn.Keychain{
		authn.DefaultKeychain,
	}
	if strings.Contains(host, "azurecr.io") {
		keychain = append(keychain, authn.NewKeychainFromHelper(acr.ACRCredHelper{}))
	} else if strings.Contains(host, "gcr.io") {
		keychain = append(keychain, google.Keychain)
	} else if strings.Contains(host, "ecr.aws") || strings.Contains(host, "amazonaws.com") {
		keychain = append(keychain, authn.NewKeychainFromHelper(ecr.NewECRHelper(ecr.WithClientFactory(api.DefaultClientFactory{}))))
	}
	return keychain
}
