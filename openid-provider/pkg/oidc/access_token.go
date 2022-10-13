package oidc

type AccessToken struct {
	ID     ID
	UserID ID
	// http form params
	Code        string
	RedirectURI string
	ClientID    string
	GrantType   string
}

func (t AccessToken) GetID() ID {
	return t.ID
}
