package controller

import (
	"net/http"

	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/storage"
)

type Token struct {
	authCodeStorage storage.AuthCode
	clientStorage   storage.Client
	tokenStorage    storage.AccessToken
}

func (t Token) SwitchCodeToToken(w http.ResponseWriter, r *http.Request) {

}

func (t Token) CheckToken(w http.ResponseWriter, r *http.Request) {

}
