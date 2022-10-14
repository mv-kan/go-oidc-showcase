package oidc

type AccessToken struct {
	// ID is access token in http json response
	ID     ID
	UserID ID
	// http form params
	Scopes      []string
	RedirectURI string
	ClientID    string
	GrantType   string
}

func (t AccessToken) GetID() ID {
	return t.ID
}
