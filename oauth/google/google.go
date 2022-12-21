package google

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/manoj-gupta/go-auth/oauth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	profileURL = "https://www.googleapis.com/oauth2/v3/userinfo?access_token="
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

// CallbackHandler validates callback, get token and returns basic user info
func (p *Provider) CallbackHandler(w http.ResponseWriter, r *http.Request) (oauth.User, error) {
	state := r.FormValue("state")
	code := r.FormValue("code")

	if state != oauth.RandomString {
		return oauth.User{}, errors.New("Invalid user state")
	}

	return p.getUserData(state, code)
}

func (p *Provider) getUserData(state, code string) (oauth.User, error) {
	user := oauth.User{
		Provider: p.Name(),
	}

	token, err := p.config.Exchange(context.Background(), code)
	if err != nil {
		return user, err
	}

	response, err := http.Get(profileURL + token.AccessToken)
	if err != nil {
		fmt.Printf("Not found user info err: %v\n", err)
		return user, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return user, err
	}

	var u googleUser
	if err := json.Unmarshal(body, &u); err != nil {
		return user, err
	}

	// Extract the user data
	user.UserID = u.ID
	user.Name = u.Name
	user.Email = u.Email
	user.AvatarURL = u.Picture

	return user, nil
}

type googleUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Link      string `json:"link"`
	Picture   string `json:"picture"`
}
