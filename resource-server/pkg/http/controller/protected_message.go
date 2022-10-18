package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/http/response"
	"github.com/mv-kan/go-oidc-showcase/openid-provider/pkg/log"
	"github.com/mv-kan/go-oidc-showcase/resource-server/pkg/config"
	"github.com/mv-kan/go-oidc-showcase/resource-server/pkg/http/utils"
)

type ProtectedMessage struct {
	opconf config.OP
	logger log.Logger
}

func NewProtectedMessage(logger log.Logger, opconf config.OP) ProtectedMessage {
	return ProtectedMessage{logger: logger, opconf: opconf}
}
func (m ProtectedMessage) GetProtectedSuperSecret(w http.ResponseWriter, r *http.Request) {
	bearerToken := r.Header.Get("Authorization")

	// get token without bearer
	tmp := strings.Split(bearerToken, "Bearer ")
	if len(tmp) != 2 {
		m.logger.Error("invalid authorization header")
		utils.WriteResponse(w, http.StatusBadRequest, "invalid authorization header")
		return
	}
	token := tmp[1]

	// check token
	// send token to server
	checkTokenURL, err := url.JoinPath(m.opconf.URL, m.opconf.CheckTokenEndpoint)
	if err != nil {
		m.logger.Error(err.Error())
		utils.WriteResponse(w, http.StatusInternalServerError, "internal error")
	}

	// get resp
	//one-line post request/resp...
	resp, err := http.PostForm(checkTokenURL, url.Values{
		"token": {token},
	})

	//okay, moving on...
	if err != nil {
		//handle postform error
		m.logger.Error(err.Error())
		utils.WriteResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	defer resp.Body.Close()
	// body, err := ioutil.ReadAll(response.Body)

	// if err != nil {
	// 	//handle read response error
	// 	m.logger.Error(err.Error())
	// 	utils.WriteResponse(w, http.StatusInternalServerError, "internal server error")
	// 	return
	// }
	if resp.StatusCode != http.StatusOK {
		m.logger.Info("access denied, access token is not valid")
		utils.WriteResponse(w, http.StatusUnauthorized, "access token is not valid")
		return
	}
	jsonResponse := response.CheckToken{}
	if err := json.NewDecoder(resp.Body).Decode(&jsonResponse); err != nil {
		m.logger.Error(err.Error())
		utils.WriteResponse(w, http.StatusInternalServerError, "internal server error")
		return
	}

	accessToken := jsonResponse.AccessToken
	userID := jsonResponse.Username

	utils.WriteResponse(w, http.StatusOK, fmt.Sprintf("Hello %s, this is a secret message and your access token is %s", userID, accessToken))
}
