package oidc

type User struct {
	// ID is username, this is simple but it is for showcase purposes
	ID       ID
	Password string
}

func (u User) GetID() ID {
	return u.ID
}
