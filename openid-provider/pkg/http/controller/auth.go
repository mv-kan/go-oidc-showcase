package controller

import (
	"net/http"

	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/storage"
)

type Auth struct {
	requestStorage storage.AuthRequest
}

func (a Auth) Auth(w http.ResponseWriter, r *http.Request) {

}
