package config

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
	Endpoint     oauth2.Endpoint
}

func NewGoogleOAuthConfig(cID, cSecret, redirectURL string, scopes []string) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     cID,
		ClientSecret: cSecret,
		RedirectURL:  redirectURL,
		Scopes:       scopes,
		Endpoint:     google.Endpoint,
	}
}
