package oidc

type ID = string
type IDer interface {
	GetID() ID
}
