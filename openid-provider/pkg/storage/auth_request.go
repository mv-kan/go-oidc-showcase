package storage

import "github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/oidc"

type AuthRequest interface {
	storage[oidc.AuthRequest]
}

func NewAuthRequest() AuthRequest {
	return newStorage[oidc.AuthRequest]()
}
