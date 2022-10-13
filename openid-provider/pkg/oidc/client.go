package oidc

type Client struct {
	ID          ID
	Secret      string
	RedirectURI []string
}

func (c Client) GetID() ID {
	return c.ID
}
