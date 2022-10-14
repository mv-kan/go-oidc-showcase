package oidc

type AuthCode struct {
	ID     ID
	UserID ID
	// http form params
	ClientID     string
	RedirectURI  string
	State        string
	Scope        []string
	ResponseType []string
}

func (c AuthCode) GetID() ID {
	return c.ID
}
