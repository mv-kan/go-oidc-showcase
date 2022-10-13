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

func (a AccessToken) GetID() ID {
	return a.ID
}
