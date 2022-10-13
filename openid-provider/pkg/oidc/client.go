package oidc

type Client struct {
	ID          string
	Secret      string
	RedirectURI []string
}

func (c Client) GetID() string {
	return c.ID
}
