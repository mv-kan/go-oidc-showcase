package utils

import (
	"net/http"
	"net/url"
	"strings"
)

func SetCookie(w http.ResponseWriter, name, value string) {
	cookie := &http.Cookie{
		Name:  name,
		Value: value,
	}
	http.SetCookie(w, cookie)
}

// ErrNotFound if cookie is not found
func GetCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// grant_type = authorization_code this is always because I didnt implmenet other grant types in openid provider
// opEndpoint is token exchange endpoint
func GetExchangeCodeToTokenRequest(opEndpoint, clientID, clientSecret, redirectURI, code, grantType string) (*http.Request, error) {
	data := url.Values{
		"code":         {code},
		"redirect_uri": {redirectURI},
		"client_id":    {clientID},
		"grant_type":   {grantType},
	}
	request, err := http.NewRequest(http.MethodPost, opEndpoint, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	request.SetBasicAuth(clientID, clientSecret)
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return request, nil
}

// creates authorization code request url from arguments
// response type and scope are openid in this project, because I did not implement other types and scopes LOL
// opEndpoint is authorization endpoint
func GetAuthCodeURL(opEndpoint, clientID, redirectURI, state, scope, responseType string) (string, error) {
	opURL, err := url.Parse(opEndpoint)
	if err != nil {
		return "", err
	}
	values := opURL.Query()
	values.Add("client_id", clientID)
	values.Add("redirect_uri", redirectURI)
	values.Add("state", state)
	values.Add("scope", scope)
	values.Add("response_type", responseType)
	opURL.RawQuery = values.Encode()
	return opURL.String(), nil
}
