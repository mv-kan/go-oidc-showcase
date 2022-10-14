package oidc

type Client struct {
	ID           ID
	Secret       string
	RedirectURIs []string
}

func (c Client) GetID() ID {
	return c.ID
}
