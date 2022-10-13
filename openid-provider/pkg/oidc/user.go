package oidc

type User struct {
	ID       string
	Password string
}

func (u User) GetID() string {
	return u.ID
}
