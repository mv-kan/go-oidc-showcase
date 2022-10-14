package controller

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/http/utils"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/log"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/oidc"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/storage"
)

var (
	loginTmpl, _ = template.New("login").Parse(`
	<!DOCTYPE html>
	<html>
		<head>
			<meta charset="UTF-8">
			<title>Login</title>
		</head>
		<body style="display: flex; align-items: center; justify-content: center; height: 100vh;">
			<form method="POST" action="/login" style="height: 200px; width: 200px;">
				<input type="hidden" name="id" value="{{.ID}}">
				<div>
					<label for="username">Username:</label>
					<input id="username" name="username" style="width: 100%">
				</div>
				<div>
					<label for="password">Password:</label>
					<input id="password" name="password" style="width: 100%">
				</div>
				<p style="color:red; min-height: 1rem;">{{.Error}}</p>
				<button type="submit">Login</button>
			</form>
		</body>
	</html>`)
)

type Login struct {
	logger          log.Logger
	userStorage     storage.User
	requestStorage  storage.AuthRequest
	authCodeStorage storage.AuthCode
}

func (l Login) Login(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	authReqIDs, ok := params["authRequestID"]
	if !ok {
		l.logger.Debug("not parsed properly")
		utils.WriteResponse(w, http.StatusInternalServerError, "")
		return
	}
	authReqID := authReqIDs[0]

	_, err := l.requestStorage.Get(authReqID)
	if errors.Is(err, storage.ErrNotFound) {
		l.logger.Debug("authenticate request does not exist")
		utils.ResponseJSON(w, http.StatusForbidden, map[string]string{"error": "auth request does not exist"})
		return
	} else if err != nil {
		l.logger.Error(err.Error())
		utils.WriteResponse(w, http.StatusInternalServerError, "")
		return
	}
	renderLogin(l.logger, w, authReqID, nil)
}

func (l Login) CheckLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Error(err.Error())
		http.Error(w, fmt.Sprintf("cannot parse form:%s", err), http.StatusInternalServerError)
		return
	}
	username := r.FormValue("username")
	password := r.FormValue("password")
	id := r.FormValue("id")
	err = utils.CheckUsernamePassword(l.userStorage, l.requestStorage, username, password, id)
	if err != nil {
		log.Error(err.Error())
		renderLogin(l.logger, w, id, err)
		return
	}
	authReq, err := l.requestStorage.Get(id)
	if err != nil {
		log.Debug(err.Error())
		renderLogin(l.logger, w, id, err)
		return
	}
	authCode, err := utils.GenerateAuthCode(l.requestStorage, l.authCodeStorage, authReq)
	authCode.UserID = username
	if err != nil {
		log.Debug(err.Error())
		renderLogin(l.logger, w, id, err)
		return
	}
	// add auth code to storage
	url, err := getCallbackURL(authCode)
	if err != nil {
		// We use error function because this is out of reach to user, it purely server side thing
		log.Error(err.Error())
		renderLogin(l.logger, w, id, err)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}

func renderLogin(logger log.Logger, w http.ResponseWriter, id string, err error) {
	var errMsg string
	if err != nil {
		errMsg = err.Error()
	}
	data := &struct {
		ID    string
		Error string
	}{
		ID:    id,
		Error: errMsg,
	}
	err = loginTmpl.Execute(w, data)
	if err != nil {
		logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func getCallbackURL(authCode oidc.AuthCode) (string, error) {
	redirectURL, err := url.Parse(authCode.RedirectURI)
	if err != nil {
		return "", err
	}

	values := redirectURL.Query()
	values.Add("code", authCode.ID)
	values.Add("state", authCode.State)
	redirectURL.RawQuery = values.Encode()
	return redirectURL.String(), nil
}
