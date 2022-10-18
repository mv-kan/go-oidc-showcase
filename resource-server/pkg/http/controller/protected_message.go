package controller

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/log"
	"github.com/mv-kan/go-oidc-showcase/resource-server/pkg/config"
	"github.com/mv-kan/go-oidc-showcase/resource-server/pkg/http/utils"
)

type ProtectedMessage struct {
	OP config.OP
}

func (m ProtectedMessage) GetProtectedSuperSecret(w http.ResponseWriter, r *http.Request) {
	bearerToken := r.Header.Get("Authorization")

	// get token without bearer
	tmp := strings.Split(bearerToken, "Bearer ")
	if len(tmp) != 2 {
		log.Error("invalid authorization header")
		utils.WriteResponse(w, http.StatusBadRequest, "invalid authorization header")
		return
	}
	token := tmp[1]

	// check token
	// send token to server
	checkTokenURL, err := url.JoinPath(m.OP.URL, m.OP.CheckTokenEndpoint)
	if err != nil {
		log.Error(err.Error())
		utils.WriteResponse(w, http.StatusInternalServerError, "internal error")
	}

	// get response
	//one-line post request/response...
	response, err := http.PostForm(checkTokenURL, url.Values{
		"token": {token},
	})

	//okay, moving on...
	if err != nil {
		//handle postform error
		log.Error(err.Error())
		utils.WriteResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	defer response.Body.Close()
	// body, err := ioutil.ReadAll(response.Body)

	// if err != nil {
	// 	//handle read response error
	// 	log.Error(err.Error())
	// 	utils.WriteResponse(w, http.StatusInternalServerError, "internal server error")
	// 	return
	// }
	if response.StatusCode != http.StatusOK {
		log.Info("access denied, access token is not valid")
		utils.WriteResponse(w, http.StatusUnauthorized, "access token is not valid")
		return
	}
}
