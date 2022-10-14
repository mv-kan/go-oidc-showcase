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

func newStorage[T oidc.IDer]() storage[T] {
	models := make([]T, 0)
	return &modelStorage[T]{models: models}
}

type modelStorage[T oidc.IDer] struct {
	models []T
}

func (s modelStorage[T]) Get(id oidc.ID) (T, error) {
	for _, model := range s.models {
		if model.GetID() == id {
			return model, nil
		}
	}
	// I cannot return T like T{}, only var tmp T
	var tmp T
	return tmp, ErrNotFound
}

func (s modelStorage[T]) GetAll() ([]T, error) {
	return s.models, nil
}

func (s modelStorage[T]) Add(model T) error {
	s.models = append(s.models, model)
	return nil
}

func (s modelStorage[T]) Remove(id oidc.ID) error {
	for index, model := range s.models {
		if model.GetID() == id {
			s.models = remove(s.models, index)
			return nil
		}
	}
	return ErrNotFound
}

func remove[T oidc.IDer](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
