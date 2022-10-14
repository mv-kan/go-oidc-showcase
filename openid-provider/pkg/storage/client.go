package storage

import "github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/oidc"

type Client interface {
	storage[oidc.Client]
}

func NewClient() Client {
	return newStorage[oidc.Client]()
}
