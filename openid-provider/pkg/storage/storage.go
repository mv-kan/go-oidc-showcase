package storage

import "github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/oidc"

// this storage is not exported because
// seperation of storages by files and manually
// is more clear than separation in code by passing types
// ErrNotFound is error for storage interface
type storage[T oidc.IDer] interface {
	Get(id oidc.ID) (T, error)
	GetAll() ([]T, error)
	Add(model T) error
	Remove(id oidc.ID) error
}
