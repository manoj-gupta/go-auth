package github

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/manoj-gupta/go-auth/oauth"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

const (
	profileURL = "https://api.github.com/user"
)

// New creates a new github provider
func New(clientKey, clientSecret, redirectURL string, scopes ...string) *Provider {
	p := &Provider{
		ClientKey:    clientKey,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		name:         "github",
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
		Endpoint:     github.Endpoint,
		Scopes:       []string{},
	}

	if len(scopes) > 0 {
		c.Scopes = append(c.Scopes, scopes...)
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

	req, err := http.NewRequest("GET", profileURL, nil)
	if err != nil {
		fmt.Printf("can't request %s: %v", profileURL, err)
		return user, err
	}
	req.Header.Add("Authorization", "Bearer "+token.AccessToken)

	response, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("Not found user info err: %v\n", err)
		return user, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return user, err
	}

	err = json.NewDecoder(bytes.NewReader(body)).Decode(&user.Any)
	if err != nil {
		return user, err
	}

	// Extract the user data
	err = userFromReader(bytes.NewReader(body), &user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func userFromReader(reader io.Reader, user *oauth.User) error {
	u := struct {
		ID       int    `json:"id"`
		Email    string `json:"email"`
		Bio      string `json:"bio"`
		Name     string `json:"name"`
		Login    string `json:"login"`
		Picture  string `json:"avatar_url"`
		Location string `json:"location"`
	}{}

	err := json.NewDecoder(reader).Decode(&u)
	if err != nil {
		return err
	}

	user.UserID = strconv.Itoa(u.ID)
	user.Email = u.Email
	user.Name = u.Name
	user.Location = u.Location
	user.AvatarURL = u.Picture
	user.Description = u.Bio

	return err
}
