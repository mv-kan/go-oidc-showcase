package oidc

type User struct {
	ID       ID
	Password string
}

func (u User) GetID() ID {
	return u.ID
}
