package oauth

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
)

var randomString = func() string {
	nonceBytes := make([]byte, 64)
	_, err := io.ReadFull(rand.Reader, nonceBytes)
	if err != nil {
		panic("random string not generated err: " + err.Error())
	}
	return base64.URLEncoding.EncodeToString(nonceBytes)
}

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

	url := provider.GetAuthURL(randomString())

	return url, nil
}
