package controllers

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var SSOObj *oauth2.Config
var RandomString = "randomstrinASDg"

func init() {

	SSOObj = &oauth2.Config{
		RedirectURL:  "https://localhost:8000/auth/google/callback",
		ClientID:     "",
		ClientSecret: "",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

func Signin(w http.ResponseWriter, r *http.Request) {
	url := SSOObj.AuthCodeURL(RandomString)
	fmt.Println(url)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
