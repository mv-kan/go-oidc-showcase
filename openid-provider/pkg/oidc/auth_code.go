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
