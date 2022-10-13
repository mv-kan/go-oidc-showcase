package oidc

type AuthCode struct {
	ID string
	// http form params
	ClientID     string
	RedirectURI  string
	State        string
	Scope        []string
	ResponseType []string
}

func (a AuthCode) GetID() string {
	return a.ID
}
