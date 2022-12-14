package storage

import "github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/oidc"

type User interface {
	storage[oidc.User]
}

func NewUser() User {
	return newStorage[oidc.User]()
}
