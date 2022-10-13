package storage

import "github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/oidc"

type AccessToken interface {
	storage[oidc.AccessToken]
}
