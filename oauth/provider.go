// Package oauth provides a common API for integrating with oauth2 providers
package oauth

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// GetProviderName returns the name of a provider from URL query string
func GetProviderName(r *http.Request) (string, error) {
	// get it from the url param "provider"
	if p := r.URL.Query().Get("provider"); p != "" {
		return p, nil
	}

	// get it from the url param "provider"
	if p := r.URL.Query().Get(":provider"); p != "" {
		return p, nil
	}

	// get it from the context's value of "provider" key
	if p, ok := mux.Vars(r)["provider"]; ok {
		return p, nil
	}

	return "", errors.New("invalid provider")
}

// Provider interface to be implement for each authenticator provier (e.g. Github)
type Provider interface {
	Name() string
	GetAuthURL(state string) string
}

// Providers is a list of registered providers
type Providers map[string]Provider

var providers = Providers{}

// RegisterProvider registers a list of providers for use
func RegisterProvider(plist ...Provider) {
	for _, provider := range plist {
		providers[provider.Name()] = provider
	}
}

// GetProviders returns a list of providers
func GetProviders() Providers {
	return providers
}

// GetProvider returns a registered provider
func GetProvider(name string) (Provider, error) {
	provider := providers[name]
	if provider == nil {
		return nil, fmt.Errorf("no provider for %s exists", name)
	}
	return provider, nil
}

// UnregisterProvider registers a list of providers for use
func UnregisterProvider(name string) {
	delete(providers, name)
}
