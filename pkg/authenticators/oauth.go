package authenticator

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Authenticator is used to authenticate our users.
type Authenticator struct {
	*oidc.Provider
	oauth2.Config
	Options []oauth2.AuthCodeOption
}

// New instantiates the *Authenticator.
func New() (*Authenticator, error) {
	issuer := fmt.Sprintf("https://%s/", os.Getenv("AUTH0_DOMAIN"))
	provider, err := oidc.NewProvider(
		context.Background(),
		issuer,
	)
	if err != nil {
		return nil, err
	}

	conf := oauth2.Config{
		ClientID:     os.Getenv("AUTH0_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH0_CLIENT_SECRET"),
		RedirectURL:  os.Getenv("AUTH0_CALLBACK_URL"),
		Endpoint:     provider.Endpoint(),
		Scopes: []string{
			oidc.ScopeOpenID,
			"profile",
			"isAdmin",
			// fmt.Sprintf("%sisAdmin", issuer),
		},
	}

	return &Authenticator{
		Provider: provider,
		Config:   conf,
		Options: []oauth2.AuthCodeOption{
			oauth2.SetAuthURLParam(
				"audience",
				os.Getenv("AUTH0_AUDIENCE"),
			),
		},
	}, nil
}

// VerifyIDToken verifies that an *oauth2.Token is a valid *oidc.IDToken.
func (a *Authenticator) VerifyIDToken(ctx context.Context, token *oauth2.Token) (*oidc.IDToken, error) {
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token field in oauth2 token")
	}

	oidcConfig := &oidc.Config{
		ClientID: a.ClientID,
	}

	return a.Verifier(oidcConfig).Verify(ctx, rawIDToken)
}
