package base

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
)

var githubOauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
	ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
	Scopes:       []string{"user"},
	Endpoint:     github.Endpoint,
	// RedirectURL:  "http://localhost:8000/auth/github/callback",
}

func githubOauthLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Login 1")
	oauthState := generateStateOauthCookie(w)

	fmt.Printf("Login 2")
	u := githubOauthConfig.AuthCodeURL(oauthState)
	fmt.Printf(u)
	http.Redirect(w, r, "/auth/github/callback", http.StatusTemporaryRedirect)
}

func githubOauthCallback(w http.ResponseWriter, r *http.Request) {
	oauthState, _ := r.Cookie("oauthstate")

	fmt.Printf("Call back 1")
	fmt.Println(r.FormValue("state"))
	if r.FormValue("state") != oauthState.Value {
		log.Println("Invalid Oauth github state")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Printf("Call back 2")
	data, err := getUserDataFromGithub(r.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	fmt.Fprintf(w, "User Info: %s\n", data)
}

func getUserDataFromGithub(code string) ([]byte, error) {
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}
