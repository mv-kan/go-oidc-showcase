package response

type CheckToken struct {
	AccessToken string `json:"access_token"`
	Username    string `json:"username"`
}
