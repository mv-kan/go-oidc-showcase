package storage

import "github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/oidc"

type storage[T oidc.IDer] interface {
	Get(id oidc.ID) T
	GetAll() []T
}
