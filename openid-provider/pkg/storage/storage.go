package storage

import "github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/oidc"

// this storage is not exported because
// seperation of storages by files and manually
// is more clear than separation in code by passing
// types
type storage[T oidc.IDer] interface {
	Get(id oidc.ID) T
	GetAll() []T
}
