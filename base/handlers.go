package base

import (
	"net/http"
)

func New() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir("public/")))

	mux.HandleFunc("/auth/google/login", googleOauthLogin)
	mux.HandleFunc("/auth/google/callback", googleOauthCallback)

	return mux
}