package config

// config op (openid provider) contains all config data for op
// needed for resource server to work
type OP struct {
	URL                string
	CheckTokenEndpoint string
}
