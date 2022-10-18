package response

type GrantToken struct {
	AccessToken string
	TokenType   string
	ExpiresIn   int
	IDToken     string
}
