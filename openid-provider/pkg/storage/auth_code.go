package storage

import "github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/oidc"

type AuthCode interface {
	storage[oidc.AuthCode]
}
