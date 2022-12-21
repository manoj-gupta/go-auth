package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

// RandomString is generated randomly at start to be used for state in AuthURL and CallbackHandler
var RandomString = func() string {
	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		panic("random string not generated err: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
}()

// SignInHandler wraps providers signin functions
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	url, err := getAuthURL(w, r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err)
		return
	}
	fmt.Printf("Got url: %s\n", url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func getAuthURL(w http.ResponseWriter, r *http.Request) (string, error) {
	name, err := GetProviderName(r)
	if err != nil {
		return "", err
	}

	provider, err := GetProvider(name)
	if err != nil {
		return "", err
	}

	url := provider.GetAuthURL(RandomString)

	return url, nil
}

// CallbackHandler wraps providers callback functions
func CallbackHandler(w http.ResponseWriter, r *http.Request) (User, error) {
	name, err := GetProviderName(r)
	if err != nil {
		return User{}, err
	}

	provider, err := GetProvider(name)
	if err != nil {
		return User{}, err
	}

	return provider.CallbackHandler(w, r)
}

// Logout wraps providers logout functions
func Logout(w http.ResponseWriter, r *http.Request) error {
	// TBD
	return nil
}
