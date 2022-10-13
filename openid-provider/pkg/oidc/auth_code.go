package oidc

type AuthCode struct {
	ID ID
	// http form params
	ClientID     string
	RedirectURI  string
	State        string
	Scope        []string
	ResponseType []string
}

func (a AuthCode) GetID() ID {
	return a.ID
}
