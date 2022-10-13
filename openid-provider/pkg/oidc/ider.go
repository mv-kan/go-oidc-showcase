package oidc

type IDer interface {
	GetID() string
}
