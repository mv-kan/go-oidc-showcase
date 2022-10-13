package controller

import (
	"net/http"

	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/storage"
)

type Login struct {
	requestStorage  storage.AuthRequest
	authCodeStorage storage.AuthCode
}

func (l Login) Login(w http.ResponseWriter, r *http.Request) {

}

func (l Login) CheckLogin(w http.ResponseWriter, r *http.Request) {

}
