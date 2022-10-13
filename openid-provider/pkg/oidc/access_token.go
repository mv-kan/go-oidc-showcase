package oidc

type AccessToken struct {
	ID     string
	UserID string
	// http form params
	Code        string
	RedirectURI string
	ClientID    string
	GrantType   string
}

func (a AccessToken) GetID() string {
	return a.ID
}
