package google

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// New creates a new google provider
func New(clientKey, clientSecret, redirectURL string, scopes ...string) *Provider {
	p := &Provider{
		ClientKey:    clientKey,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		name:         "google",
	}
	p.config = newConfig(p, scopes)

	return p
}

// Provider implements oauth.Provider interface
type Provider struct {
	ClientKey    string
	ClientSecret string
	RedirectURL  string
	config       *oauth2.Config
	name         string
}

// Name is the name of this provider
func (p *Provider) Name() string {
	return p.name
}

// GetAuthURL returns redirection URL
func (p *Provider) GetAuthURL(state string) string {
	return p.config.AuthCodeURL(state)
}

func newConfig(p *Provider, scopes []string) *oauth2.Config {
	c := &oauth2.Config{
		ClientID:     p.ClientKey,
		ClientSecret: p.ClientSecret,
		RedirectURL:  p.RedirectURL,
		Endpoint:     google.Endpoint,
		Scopes:       []string{},
	}

	if len(scopes) > 0 {
		c.Scopes = append(c.Scopes, scopes...)
	} else {
		c.Scopes = []string{"email"} // fallback scope
	}

	return c
}
